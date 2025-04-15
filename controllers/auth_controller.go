package controllers

import (
	"net/http"

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

    user, err := services.RegisterUser(req)
    utils.PanicIfError(err)

    webResponse := utils.WebResponseSuccess{
        Success: true,
        Code:    http.StatusCreated,
        Status:  "CREATED",
        Message: "User created successfully, please check your email to verify your account",
        Data:    UserData{
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

    err = services.VerifyUserEmail(req.Email, req.Code)
    utils.PanicIfError(err)

    webResponse := utils.WebResponseSuccess{
        Success: true,
        Code:    http.StatusOK,
        Status:  "OK",
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

	err = services.ResendVerificationCode(req.Email)
	utils.PanicIfError(err)

	webResponse := utils.WebResponseSuccess{
		Success: true,
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Verification code resent successfully, please check your email",
	}

	w.WriteHeader(http.StatusOK)
	utils.WriteToResponseBody(w, webResponse)
}


func Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    var req validations.LoginRequest
    utils.ReadFromRequestBody(r, &req)
    err := validate.Struct(req)
    utils.PanicIfError(err)

    tokens, err := services.AuthenticateUser(req.Email, req.Password)
    utils.PanicIfError(err)

    webResponse := utils.WebResponseSuccess{
        Success: true,
        Code:    http.StatusOK,
        Status:  "OK",
        Message: "Login success",
        Data:    tokens,
    }

    w.WriteHeader(http.StatusOK)
    utils.WriteToResponseBody(w, webResponse)
}


func RefreshToken(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req validations.RefreshTokenRequest
	utils.ReadFromRequestBody(r, &req)

	err := validate.Struct(req)
	utils.PanicIfError(err)

	accessToken, err := services.RefreshToken(req.RefreshToken)
	utils.PanicIfError(err)

	webResponse := utils.WebResponseSuccess{
		Success: true,
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Token refreshed successfully",
		Data: map[string]string{
			"access_token": accessToken["access_token"],
		},
	}

	w.WriteHeader(http.StatusOK)
	utils.WriteToResponseBody(w, webResponse)
}


func GetProfile(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	id := r.Header.Get("X-User-ID")

	var user models.User
	err := config.DB.First(&user, id).Error
	if err != nil {
		panic(errors.CustomError{
			Code: http.StatusNotFound,
			Status: "NOT FOUND",
			Message: "User not found",
		})
	}

	webResponse := utils.WebResponseSuccess{
		Success: true,
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "User profile fetched successfully",
		Data: map[string]interface{}{
			"id":       user.ID,
			"name":     user.Name,
			"email":    user.Email,
			"verified": user.Verified,
		},
	}

	w.WriteHeader(http.StatusOK)
	utils.WriteToResponseBody(w, webResponse)
}

