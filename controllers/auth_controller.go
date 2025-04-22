package controllers

import (
	"net/http"

	"github.com/abdultalif/golang-auth-gorm/config"
	"github.com/abdultalif/golang-auth-gorm/errors"
	"github.com/abdultalif/golang-auth-gorm/logger"
	"github.com/abdultalif/golang-auth-gorm/models"
	"github.com/abdultalif/golang-auth-gorm/services"
	"github.com/abdultalif/golang-auth-gorm/utils"
	"github.com/abdultalif/golang-auth-gorm/validations"
	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

var validate = validator.New()
type UserData struct {
    ID       uint   `json:"id"`
    Name     string `json:"name"`
    Email    string `json:"email"`
    Verified bool   `json:"verified"`
}

func Register(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    logger.Log.WithFields(logrus.Fields{
        "method": r.Method,
        "path":   r.URL.Path,
    }).Info("Register endpoint called")

    var req validations.RegisterRequest
    utils.ReadFromRequestBody(r, &req)
    logger.Log.WithFields(logrus.Fields{
        "email": req.Email,
        "name":  req.Name,
    }).Info("Registration attempt")

    err := validate.Struct(req)
    if err != nil {
        logger.Log.WithFields(logrus.Fields{
            "error": err.Error(),
        }).Error("Validation failed")
        panic(err)
    }

    user, err := services.RegisterUser(req)
    if err != nil {
        logger.Log.WithFields(logrus.Fields{
            "error": err.Error(),
            "email": req.Email,
        }).Error("Registration failed")
        panic(err)
    }

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

    logger.Log.WithFields(logrus.Fields{
        "id": user.ID,
        "email":   user.Email,
    }).Info("User registered successfully and sent verification email")

    w.WriteHeader(http.StatusCreated)
    utils.WriteToResponseBody(w, webResponse)
}

func VerifyOTP(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

    logger.Log.WithFields(logrus.Fields{
        "method": r.Method,
        "path":   r.URL.Path,
    }).Info("Verification OTP endpoint called")

    var req validations.VerifyCodeRequest
    utils.ReadFromRequestBody(r, &req)

    logger.Log.WithFields(logrus.Fields{
        "email": req.Email,
        "code":  req.Code,
    }).Info("Verification OTP attempt")
    


    err := validate.Struct(req)
	if err != nil {
        logger.Log.WithFields(logrus.Fields{
            "error": err.Error(),
        }).Error("Validation failed")
        panic(err)
    }
    

    err = services.VerifyUserEmail(req.Email, req.Code)
    if err != nil {
        logger.Log.WithFields(logrus.Fields{
            "error": err.Error(),
        }).Error("Verification OTP failed")
        panic(err)
    }

    webResponse := utils.WebResponseSuccess{
        Success: true,
        Code:    http.StatusOK,
        Status:  "OK",
        Message: "User verified successfully",
    }

    logger.Log.Info("User verified successfully")

    w.WriteHeader(http.StatusOK)
    utils.WriteToResponseBody(w, webResponse)
}

func ResendCode(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

    logger.Log.WithFields(logrus.Fields{
        "method": r.Method,
        "path":   r.URL.Path,
    }).Info("Resend code endpoint called")

	var req validations.ResendCodeRequest
	utils.ReadFromRequestBody(r, &req)

    logger.Log.WithFields(logrus.Fields{
        "email": req.Email,
    }).Info("Resend code attempt")

	err := validate.Struct(req)
	if err != nil {
        logger.Log.WithFields(logrus.Fields{
            "error": err.Error(),
        }).Error("Validation failed")
        panic(err)
    }

	err = services.ResendVerificationCode(req.Email)
	if err != nil {
        logger.Log.WithFields(logrus.Fields{
            "error": err.Error(),
        }).Error("Resend code failed")
        panic(err)
    }

	webResponse := utils.WebResponseSuccess{
		Success: true,
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Verification code resent successfully, please check your email",
	}

    logger.Log.Info("Verification code resent successfully")

	w.WriteHeader(http.StatusOK)
	utils.WriteToResponseBody(w, webResponse)
}


func Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    logger.Log.WithFields(logrus.Fields{
        "method": r.Method,
        "path":   r.URL.Path,
    }).Info("Login endpoint called")

    var req validations.LoginRequest
    utils.ReadFromRequestBody(r, &req)
    logger.Log.WithFields(logrus.Fields{
        "email": req.Email,
    }).Info("Loginn attempt")


    err := validate.Struct(req)
    if err != nil {
        logger.Log.WithFields(logrus.Fields{
            "error": err.Error(),
        }).Error("Validation failed")
        panic(err)  
    }

    tokens, err := services.AuthenticateUser(req.Email, req.Password)
    if err != nil {
        logger.Log.WithFields(logrus.Fields{
            "error": err.Error(),
        }).Error("Login failed")
        panic(err)
    }

    webResponse := utils.WebResponseSuccess{
        Success: true,
        Code:    http.StatusOK,
        Status:  "OK",
        Message: "Login success",
        Data:    tokens,
    }

    logger.Log.WithFields(logrus.Fields{
        "access_token": "secret",
        "refresh_token": "secret",
    }).Info("User login successfully")

    w.WriteHeader(http.StatusOK)
    utils.WriteToResponseBody(w, webResponse)
}


func RefreshToken(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    logger.Log.WithFields(logrus.Fields{
        "method": r.Method,
        "path":   r.URL.Path,
    }).Info("Refresh token endpoint called")

	var req validations.RefreshTokenRequest
	utils.ReadFromRequestBody(r, &req)
    logger.Log.WithFields(logrus.Fields{
        "refresh_token": "secret",
    }).Info("Refresh token attempt")

	err := validate.Struct(req)
    if err != nil {
        logger.Log.WithFields(logrus.Fields{
            "error": err.Error(),
        }).Error("Validation failed")
        panic(err)
    }

	accessToken, err := services.RefreshToken(req.RefreshToken)
	if err != nil {
        logger.Log.WithFields(logrus.Fields{
            "error": err.Error(),
        }).Error("Refresh token failed")
        panic(err)
    }

	webResponse := utils.WebResponseSuccess{
		Success: true,
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Token refreshed successfully",
		Data: map[string]string{
			"access_token": accessToken["access_token"],
		},
	}

    logger.Log.WithFields(logrus.Fields{
        "access_token": "secret",
    }).Info("User refresh token successfully")

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

    logger.Log.WithFields(logrus.Fields{
        "access_token": "secret",
    }).Info("User refresh token successfully")

	w.WriteHeader(http.StatusOK)
	utils.WriteToResponseBody(w, webResponse)
}


func ForgotPassword(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    
    logger.Log.WithFields(logrus.Fields{
        "method": r.Method,
        "path":   r.URL.Path,
    }).Info("Forgot password endpoint called")

    var req validations.ForgotPasswordRequest
    utils.ReadFromRequestBody(r, &req)
    logger.Log.WithFields(logrus.Fields{
        "email": req.Email,
    }).Info("Forgot password attempt")

    err := validate.Struct(req)
    if err != nil {
        logger.Log.WithFields(logrus.Fields{
            "error": err.Error(),
        }).Error("Validation failed")
        panic(err)
    }

    err = services.ForgotPassword(req.Email)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Forgot password service failed")
		panic(err)
	}

    webResponse := utils.WebResponseSuccess{
		Success: true,
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Forgot password successfully",
	}

    logger.Log.Info("Forgot password successfully")

	w.WriteHeader(http.StatusOK)
	utils.WriteToResponseBody(w, webResponse)

}

func ResetPassword(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
}
func VerifyResetOTP(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
}
