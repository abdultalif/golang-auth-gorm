package services

import (
	"net/http"
	"time"

	"github.com/abdultalif/golang-auth-gorm/config"
	"github.com/abdultalif/golang-auth-gorm/errors"
	"github.com/abdultalif/golang-auth-gorm/logger"
	"github.com/abdultalif/golang-auth-gorm/models"
	"github.com/abdultalif/golang-auth-gorm/utils"
	"github.com/abdultalif/golang-auth-gorm/validations"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)


func RegisterUser(req validations.RegisterRequest) (*models.User, error) {
 
    logger.Log.Info("Checking email existence")
    if config.DB.Where("email = ?", req.Email).First(&models.User{}).RowsAffected > 0 {
        logger.Log.WithFields(logrus.Fields{
            "email": req.Email,
        }).Warn("Email already exists")
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
    
    logger.Log.Info("Creating new user in database")
    err := config.DB.Create(&user).Error
    if err != nil {
        logger.Log.WithFields(logrus.Fields{
            "error": err.Error(),
        }).Error("Failed to create user")
        return nil, err
    }

    logger.Log.Info("Creating verification code")
    code, err := CreateVerificationCode(user.ID)
    if err != nil {
        logger.Log.WithFields(logrus.Fields{
            "error": err.Error(),
        }).Error("Failed to create verification code")
        return nil, err
    }

    logger.Log.Info("Sending verification email")
    err = SendVerificationEmail(user.Email, code)
    if err != nil {
        logger.Log.WithFields(logrus.Fields{
            "error": err.Error(),
        }).Error("Failed to send verification email")
        return nil, err
    }

    logger.Log.Info("User registered successfully and verification email sent")
    return &user, nil
}

func CreateVerificationCode(userID uint) (string, error) {

    logger.Log.Info("Generating verification code")
	code, err := utils.GenerateVerificationCode()
    if err != nil {
        logger.Log.WithFields(logrus.Fields{            
            "error":   err.Error(),
        }).Error("Failed to generate verification code")
        return "", err
    }

	hashed, _ := utils.HashPassword(code)
	verif := models.VerificationCode{
		UserID:    userID,
		Code:      hashed,
		ExpiresAt: time.Now().Add(30 * time.Minute),
	}

    logger.Log.Info("Creating verification code in database")
	err = config.DB.Create(&verif).Error

    if err != nil {
        logger.Log.WithFields(logrus.Fields{
            "error":   err.Error(),
        }).Error("Failed to create verification code in database")
        return "", err
    }

    logger.Log.Info("Verification code created successfully")

	return code, err
}

func ForgotPassword(email string) error {  

    logger.Log.WithFields(logrus.Fields{
		"email": email,
	}).Info("Checking if user exists for forgot password")

    logger.Log.Info("Checking if user is registered")
	var user models.User
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
        logger.Log.WithFields(logrus.Fields{
            "error": err.Error(),
            }).Warn("Email not registered")
            return errors.CustomError{
                Code:    http.StatusNotFound,
                Status:  "NOT FOUND",
                Message: "Email not registered",
            }
        }
        
    logger.Log.Info("Checking if user is verified")
    if !user.Verified {
        logger.Log.WithFields(logrus.Fields{
            "email": email,
        }).Warn("User not verified")
        return errors.CustomError{
            Code:    http.StatusUnauthorized,
            Status:  "UNAUTHORIZED",
            Message: "User not verified",
        }
    }

    logger.Log.Info("Deleting previous password reset token if any")
    if err := config.DB.Unscoped().Where("user_id = ?", user.ID).Delete(&models.PasswordReset{}).Error; err != nil {
        logger.Log.WithFields(logrus.Fields{
            "error": err.Error(),
        }).Error("Failed to delete old reset tokens")
        panic(err)
    }

	token := uuid.New().String()
	hashed, _ := utils.HashPassword(token)

    reset := models.PasswordReset{
		UserID:    user.ID,
		Token:      hashed,
		ExpiresAt: time.Now().Add(30 * time.Minute),
	}

	logger.Log.Info("Saving reset code to database")
	if err := config.DB.Create(&reset).Error; err != nil {
		logger.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to store reset code")
        panic(err)
	}

	
	logger.Log.Info("Sending reset password email")
	if err := SendResetPasswordEmail(email, token); err != nil {
		logger.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to send reset password email")
        panic(err)
	}

	logger.Log.Info("Reset password email sent successfully")
	return nil
}

func VerifyUserEmail(email, code string) error {

    logger.Log.WithFields(logrus.Fields{
        "email": email,
    }).Info("Attempting to verify user email")

    var user models.User
    err := config.DB.Where("email = ?", email).First(&user).Error
    if err != nil {
        logger.Log.WithFields(logrus.Fields{
            "email": email,
            "error": err.Error(),
        }).Warn("User not found for email verification")
        return errors.CustomError{
            Code:    http.StatusNotFound,
            Status:  "NOT FOUND",
            Message: "User not found",
        }
    }

    logger.Log.Info("Looking up verification code")
    var verificationCode models.VerificationCode
    err = config.DB.Where("user_id = ?", user.ID).First(&verificationCode).Error
    if err != nil {
        logger.Log.WithFields(logrus.Fields{
            "error":   err.Error(),
        }).Warn("Verification code not found")
        return errors.CustomError{
            Code:    http.StatusBadRequest,
            Status:  "BAD REQUEST",
            Message: "Verification code not found",
        }
    }

    if !utils.CheckPasswordHash(code, verificationCode.Code) {
        logger.Log.Warn("Invalid verification code provided")
        return errors.CustomError{
            Code:    http.StatusBadRequest,
            Status:  "BAD REQUEST",
            Message: "Invalid code",
        }
    }

    if time.Now().After(verificationCode.ExpiresAt) {
        logger.Log.Warn("Expired verification code provided")
        return errors.CustomError{
            Code:    http.StatusBadRequest,
            Status:  "BAD REQUEST",
            Message: "Code has expired",
        }
    }

    
    logger.Log.Info("Marking user as verified")
    user.Verified = true
    config.DB.Save(&user)
    config.DB.Delete(&verificationCode)

    logger.Log.Info("User email verified successfully")

    return nil
}

func ResendVerificationCode(email string) error {

    logger.Log.WithFields(logrus.Fields{
        "email": email,
    }).Info("Attempting to resend verification code")

	var user models.User
	err := config.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
        logger.Log.WithFields(logrus.Fields{
            "error": err.Error(),
        }).Warn("User not found for resending verification code")
        return errors.CustomError{
            Code:    http.StatusNotFound,
            Status:  "NOT FOUND",
            Message: "User not registered",
        }
    }

	if user.Verified {
        logger.Log.Warn("Attempt to resend verification code for already verified user")
        return errors.CustomError{
            Code:    http.StatusConflict,
            Status:  "CONFLICT",
            Message: "User already verified. No further action is needed.",
        }
    }

    logger.Log.Info("Creating new verification code")
	code, err := CreateVerificationCode(user.ID)
    if err != nil {
        logger.Log.WithFields(logrus.Fields{
            "error":   err.Error(),
        }).Error("Failed to create verification code")
        panic(err)
    }

    logger.Log.Info("Sending verification email")
	err = SendVerificationEmail(user.Email, code)
	if err != nil {
        logger.Log.WithFields(logrus.Fields{
            "error":   err.Error(),
        }).Error("Failed to send verification email")
        panic(err)
    }

    logger.Log.Info("Verification code resent successfully")

	return nil
}

func AuthenticateUser(email, password string) (map[string]string, error) {
    
    logger.Log.WithFields(logrus.Fields{
        "email": email,
    }).Info("Attempting to authenticate user")

    var user models.User
    err := config.DB.Where("email = ?", email).First(&user).Error

    if err != nil {
        logger.Log.WithFields(logrus.Fields{
            "error": err.Error(),
        }).Warn("Authentication failed - email not found")
        
        return nil, errors.CustomError{
            Code:    http.StatusUnauthorized,
            Status:  "UNAUTHORIZED",
            Message: "Email or password is incorrect",
        }
    }

    if !utils.CheckPasswordHash(password, user.Password) {
        logger.Log.Warn("Authentication failed - incorrect password")

        return nil, errors.CustomError{
            Code:    http.StatusUnauthorized,
            Status:  "UNAUTHORIZED",
            Message: "Email or password is incorrect",
        }
    }

    if !user.Verified {
        logger.Log.Warn("Authentication failed - account not verified")
        
        return nil, errors.CustomError{
            Code:    http.StatusForbidden,
            Status:  "FORBIDDEN",
            Message: "Please verify your account first",
        }
    }

    logger.Log.Info("Generating JWT tokens")

    accessToken, err := utils.GenerateJWT(user.ID, user.Name, user.Email, false)
    if err != nil {
        logger.Log.WithFields(logrus.Fields{
            "error":   err.Error(),
        }).Error("Failed to generate access token")
        
        return nil, err
    }

    refreshToken, err := utils.GenerateJWT(user.ID, user.Name, user.Email, true)
    if err != nil {
        logger.Log.WithFields(logrus.Fields{
            "error":   err.Error(),
        }).Error("Failed to generate refresh token")
        
        return nil, err
    }

    logger.Log.Info("User authenticated successfully")

    return map[string]string{
        "access_token":  accessToken,
        "refresh_token": refreshToken,
    }, nil
}

func RefreshToken(refreshToken string) (map[string]string, error) {
    logger.Log.Info("Attempting to refresh token")

    claims, err := utils.VerifyToken(refreshToken, true)
    if err != nil {
        logger.Log.WithFields(logrus.Fields{
            "error": err.Error(),
        }).Warn("Invalid refresh token provided")
        
        return nil, errors.CustomError{
            Code:    http.StatusUnauthorized,
            Status:  "UNAUTHORIZED",
            Message: "Invalid refresh token",
        }
    }

    id := uint(claims["id"].(float64))
    logger.Log.Info("Generating new access token")

    accessToken, err := utils.GenerateJWT(id, claims["name"].(string), claims["email"].(string), false)
    if err != nil {
        logger.Log.WithFields(logrus.Fields{
            "error":   err.Error(),
        }).Error("Failed to generate new access token")
        
        return nil, err
    }    

    logger.Log.Info("Token refreshed successfully")

    return map[string]string{    
        "access_token": accessToken,
    }, nil        
}