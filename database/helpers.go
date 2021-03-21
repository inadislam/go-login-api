package database

import (
	"errors"

	"github.com/inadislam/go-login-api/middlewares"
	"github.com/inadislam/go-login-api/models"
	"github.com/inadislam/go-login-api/utils"
)

func SignupHelper(user models.User) (models.User, error) {
	db := Connect()
	defer db.Close()

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return models.User{}, err
	}
	user.Password = string(hashedPassword)
	err = db.Debug().Model(&models.User{}).Create(&user).Error
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func LoginHelper(username, password string) (string, error) {
	user := UserByName(username)
	if user.ID == 0 {
		return "", errors.New("User not Found!!")
	}
	err := utils.ComparePass(user.Password, password)
	if err != nil {
		return "", err
	}
	token, err := middlewares.JwtToken(user)
	if err != nil {
		return "", err
	}
	return token, err
}

func UserByName(username string) models.User {
	db := Connect()
	defer db.Close()
	var user models.User
	err := db.Debug().Model(&models.User{}).Where("username = ?", username).Find(&user).Error
	utils.CheckErr(err)
	return user
}
