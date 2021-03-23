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

func UserById(userid uint32) (models.User, error) {
	db := Connect()
	defer db.Close()
	var user models.User
	err := db.Debug().Model(&models.User{}).Where("ID = ?", userid).Select("id, username, email, updated_at, created_at").Find(&user).Error
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func UserByEmail(useremail string) (models.User, error) {
	db := Connect()
	defer db.Close()
	var user models.User
	err := db.Debug().Model(&models.User{}).Where("Email = ?", useremail).Select("id, username, email, updated_at, created_at").Find(&user).Error
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func GetAllUser() ([]models.User, error) {
	db := Connect()
	defer db.Close()
	var users []models.User
	err := db.Debug().Order("id asc").Select("id, username, email, updated_at, created_at").Find(&users).Error
	if err != nil {
		return []models.User{}, err
	}
	return users, nil
}

func DeleteUser(userid uint32) (int64, error) {
	db := Connect()
	defer db.Close()

	del := db.Debug().Model(&models.User{}).Where("id = ?", userid).Take(&models.User{}).Delete(&models.User{})
	if del.Error != nil {
		return 0, del.Error
	}
	return del.RowsAffected, nil
}

func UserUpdate(userid uint32) (models.User, error) {
	db := Connect()
	defer db.Close()
	var user models.User
	db = db.Model(&models.User{}).Where("ID = ?", userid).Updates(
		models.User{
			Name:      user.Name,
			Username:  user.Username,
			Email:     user.Email,
			Password:  user.Password,
			UpdatedAt: user.UpdatedAt,
		},
	)
	if db.Error != nil {
		return models.User{}, db.Error
	}
	user, err := UserById(userid)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}
