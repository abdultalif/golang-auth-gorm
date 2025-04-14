package services

import (
	"time"

	"github.com/abdultalif/golang-auth-gorm/config"
	"github.com/abdultalif/golang-auth-gorm/errors"
	"github.com/abdultalif/golang-auth-gorm/models"
	"github.com/abdultalif/golang-auth-gorm/utils"
)

func CreateUser(name, email, password string) (*models.User, error) {

	if config.DB.Where("email = ?", email).First(&models.User{}).RowsAffected > 0 {
		panic(errors.NewDuplicateError("Email already exists"))
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
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}
	err = config.DB.Create(&verif).Error
	return code, err
}
