package controllers

import (
	"net/http"
	"time"

	"github.com/abdultalif/golang-auth-gorm/config"
	"github.com/abdultalif/golang-auth-gorm/errors"
	"github.com/abdultalif/golang-auth-gorm/models"
	"github.com/abdultalif/golang-auth-gorm/services"
	"github.com/abdultalif/golang-auth-gorm/utils"
	"github.com/abdultalif/golang-auth-gorm/validations"
	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
)

var validate = validator.New()
type UserData struct {
    ID       uint   `json:"id"`
    Name     string `json:"name"`
    Email    string `json:"email"`
    Verified bool   `json:"verified"`
}

func Register(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req validations.RegisterRequest
	utils.ReadFromRequestBody(r, &req)

	err := validate.Struct(req)
	utils.PanicIfError(err)

	user, err := services.CreateUser(req.Name, req.Email, req.Password)
	utils.PanicIfError(err)

	code, err := services.CreateVerificationCode(user.ID)
	utils.PanicIfError(err)

	newErr := services.SendVerificationEmail(user.Email, code)
	utils.PanicIfError(newErr)

	webResponse := utils.WebResponseSuccess{
		Success: true,
		Code: http.StatusCreated,
		Status: "CREATED",
		Message: "User created successfully, please check your email to verify your account",
		Data: UserData{
			ID:       user.ID,
			Name:     user.Name,
			Email:    user.Email,
			Verified: user.Verified,
		},
	}

	w.WriteHeader(http.StatusCreated)
	utils.WriteToResponseBody(w, webResponse)
}

func VerifyUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req validations.VerifyCodeRequest
	utils.ReadFromRequestBody(r, &req)

	err := validate.Struct(req)
	utils.PanicIfError(err)

	var user models.User
	err = config.DB.Where("email = ?", req.Email).First(&user).Error;
	if err != nil {
		panic(errors.CustomError{
			Code: http.StatusNotFound,
			Status: "NOT FOUND",
			Message: "User not found",
		})
	}


	var code models.VerificationCode
	err = config.DB.Where("user_id = ?", user.ID).First(&code).Error
	if err != nil {
		panic(errors.CustomError{
			Code: http.StatusBadRequest,
			Status: "BAD REQUEST",
			Message: "Verification code not found",
		})
	}

	if !utils.CheckPasswordHash(req.Code, code.Code) {
		panic(errors.CustomError{
			Code:    http.StatusBadRequest,
			Status:  "BAD REQUEST",
			Message: "Invalid code",
		})
	}

	if time.Now().After(code.ExpiresAt) {
		panic(errors.CustomError{
			Code:    http.StatusBadRequest,
			Status:  "BAD REQUEST",
			Message: "Code has expired",
		})
	}


	user.Verified = true
	config.DB.Save(&user)
	config.DB.Delete(&code)

	webResponse := utils.WebResponseSuccess{
		Success: true,
		Code: http.StatusOK,
		Status: "OK",
		Message: "User verified successfully",
	}

	w.WriteHeader(http.StatusOK)
	utils.WriteToResponseBody(w, webResponse)
}


func ResendCode(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req validations.ResendCodeRequest
	utils.ReadFromRequestBody(r, &req)

	err := validate.Struct(req)
	utils.PanicIfError(err)

	var user models.User
	err = config.DB.Where("email = ?", req.Email).First(&user).Error
	if err != nil {
		panic(errors.CustomError{
			Code:    http.StatusNotFound,
			Status:  "NOT FOUND",
			Message: "User not found",
		})
	}

	if user.Verified {
		panic(errors.CustomError{
			Code:    http.StatusConflict,
			Status:  "CONFLICT",
			Message: "User already verified. No further action is needed.",
		})
	}

	
	config.DB.Where("user_id = ?", user.ID).Delete(&models.VerificationCode{})

	
	code, err := services.CreateVerificationCode(user.ID)
	utils.PanicIfError(err)

	newErr := services.SendVerificationEmail(user.Email, code)
	utils.PanicIfError(newErr)

	webResponse := utils.WebResponseSuccess{
		Success: true,
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Verification code resent successfully, please check your email",
	}

	w.WriteHeader(http.StatusOK)
	utils.WriteToResponseBody(w, webResponse)
}

