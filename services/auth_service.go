package services

import (
	"net/http"
	"time"

	"github.com/abdultalif/golang-auth-gorm/config"
	"github.com/abdultalif/golang-auth-gorm/errors"
	"github.com/abdultalif/golang-auth-gorm/models"
	"github.com/abdultalif/golang-auth-gorm/utils"
	"github.com/abdultalif/golang-auth-gorm/validations"
)


func RegisterUser(req validations.RegisterRequest) (*models.User, error) {
    if config.DB.Where("email = ?", req.Email).First(&models.User{}).RowsAffected > 0 {
        return nil, errors.CustomError{
            Code:    http.StatusConflict,
            Status:  "CONFLICT",
            Message: "Email already exists",
        }
    }

    hashed, _ := utils.HashPassword(req.Password)
    user := models.User{
        Name:     req.Name,
        Email:    req.Email,
        Password: hashed,
    }
    err := config.DB.Create(&user).Error
    if err != nil {
        return nil, err
    }

    code, err := CreateVerificationCode(user.ID)
    if err != nil {
        return nil, err
    }

    err = SendVerificationEmail(user.Email, code)
    if err != nil {
        return nil, err
    }

    return &user, nil
}

func CreateUser(name, email, password string) (*models.User, error) {

	if config.DB.Where("email = ?", email).First(&models.User{}).RowsAffected > 0 {
		panic(errors.CustomError{
			Code: http.StatusConflict,
			Status: "CONFLICT",
			Message: "Email already exists",
		})
	}

	hashed, _ := utils.HashPassword(password)
	user := models.User{
		Name:     name,
		Email:    email,
		Password: hashed,
	}
	err := config.DB.Create(&user).Error
	return &user, err
}

func CreateVerificationCode(userID uint) (string, error) {
	code, err := utils.GenerateVerificationCode()
	hashed, _ := utils.HashPassword(code)
	if err != nil {
		return "", err
	}
	verif := models.VerificationCode{
		UserID:    userID,
		Code:      hashed,
		ExpiresAt: time.Now().Add(30 * time.Minute),
	}
	err = config.DB.Create(&verif).Error
	return code, err
}

func VerifyUserEmail(email, code string) error {
    var user models.User
    err := config.DB.Where("email = ?", email).First(&user).Error
    if err != nil {
        return errors.CustomError{
            Code:    http.StatusNotFound,
            Status:  "NOT FOUND",
            Message: "User not found",
        }
    }

    var verificationCode models.VerificationCode
    err = config.DB.Where("user_id = ?", user.ID).First(&verificationCode).Error
    if err != nil {
        return errors.CustomError{
            Code:    http.StatusBadRequest,
            Status:  "BAD REQUEST",
            Message: "Verification code not found",
        }
    }

    if !utils.CheckPasswordHash(code, verificationCode.Code) {
        return errors.CustomError{
            Code:    http.StatusBadRequest,
            Status:  "BAD REQUEST",
            Message: "Invalid code",
        }
    }

    if time.Now().After(verificationCode.ExpiresAt) {
        return errors.CustomError{
            Code:    http.StatusBadRequest,
            Status:  "BAD REQUEST",
            Message: "Code has expired",
        }
    }

    user.Verified = true
    config.DB.Save(&user)
    config.DB.Delete(&verificationCode)

    return nil
}

func ResendVerificationCode(email string) error {
	var user models.User
	err := config.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return errors.CustomError{
			Code:    http.StatusNotFound,
			Status:  "NOT FOUND",
			Message: "User not registered",
		}
	}

	if user.Verified {
		return errors.CustomError{
			Code:    http.StatusConflict,
			Status:  "CONFLICT",
			Message: "User already verified. No further action is needed.",
		}
	}	

	code, err := CreateVerificationCode(user.ID)
	utils.PanicIfError(err)

	err = SendVerificationEmail(user.Email, code)
	utils.PanicIfError(err)

	return nil
}

func AuthenticateUser(email, password string) (map[string]string, error) {
    var user models.User
    err := config.DB.Where("email = ?", email).First(&user).Error
    if err != nil || !utils.CheckPasswordHash(password, user.Password) {
        return nil, errors.CustomError{
            Code:    http.StatusUnauthorized,
            Status:  "UNAUTHORIZED",
            Message: "Email or password is incorrect",
        }
    }

    if !user.Verified {
        return nil, errors.CustomError{
            Code:    http.StatusForbidden,
            Status:  "FORBIDDEN",
            Message: "Please verify your account first",
        }
    }

    accessToken, err := utils.GenerateJWT(user.ID, user.Name, user.Email, false)
    if err != nil {
        return nil, err
    }

    refreshToken, err := utils.GenerateJWT(user.ID, user.Name, user.Email, true)
    if err != nil {
        return nil, err
    }

    return map[string]string{
        "access_token":  accessToken,
        "refresh_token": refreshToken,
    }, nil
}

func RefreshToken(refreshToken string) (map[string]string, error) {
	claims, err := utils.VerifyToken(refreshToken, true)
	if err != nil {
		return nil, errors.CustomError{
			Code:    http.StatusUnauthorized,
			Status:  "UNAUTHORIZED",
			Message: "Invalid refresh token",
		}
	}

	id := uint(claims["id"].(float64))

	accessToken, err := utils.GenerateJWT(id, claims["name"].(string), claims["email"].(string), false)

	if err != nil {
		return nil, err
	}	

	return map[string]string{	
		"access_token": accessToken,
	}, nil		
}
