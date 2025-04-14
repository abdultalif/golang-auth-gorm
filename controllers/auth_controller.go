package controllers

import (
	"net/http"

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
