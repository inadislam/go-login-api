package database

import (
	"errors"

	"github.com/inadislam/go-login-api/auth"
	"github.com/inadislam/go-login-api/models"
)

func SignupHelper(user models.User) (models.User, error) {
	db := Connect()
	defer db.Close()

	hashedPassword, err := auth.HashPassword(user.Password)
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

func LoginHelper(username, password string) (string, models.User, error) {
	user, err := UserByName(username)
	if user.ID == 0 && err != nil {
		return "", models.User{}, errors.New("user not found")
	}
	err = auth.ComparePass(user.Password, password)
	if err != nil {
		return "", models.User{}, err
	}
	token, err := auth.JwtToken(user)
	if err != nil {
		return "", models.User{}, err
	}
	return token, user, nil
}

func UserByName(username string) (models.User, error) {
	db := Connect()
	defer db.Close()
	var user models.User
	err := db.Debug().Model(&models.User{}).Where("username = ?", username).Find(&user).Error
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func UserById(userid uint32) (models.User, error) {
	db := Connect()
	defer db.Close()
	var user models.User
	err := db.Debug().Model(&models.User{}).Where("ID = ?", userid).Select("id, name, username, email, updated_at, created_at").Find(&user).Error
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func UserByEmail(useremail string) (models.User, error) {
	db := Connect()
	defer db.Close()
	var user models.User
	err := db.Debug().Model(&models.User{}).Where("Email = ?", useremail).Select("id, name, username, email, updated_at, created_at").Find(&user).Error
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func GetAllUser() ([]models.User, error) {
	db := Connect()
	defer db.Close()
	var users []models.User
	err := db.Debug().Order("id asc").Select("id, name, username, email, updated_at, created_at").Find(&users).Error
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

func UserUpdate(user models.User, userid uint32) (models.User, error) {
	db := Connect()
	defer db.Close()
	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		return models.User{}, err
	}
	user.Password = string(hashedPassword)
	db = db.Model(&models.User{}).Where("ID = ?", userid).Updates(
		map[string]interface{}{
			"name":       user.Name,
			"username":   user.Username,
			"email":      user.Email,
			"password":   user.Password,
			"updated_at": user.UpdatedAt,
		},
	)
	if db.Error != nil {
		return models.User{}, db.Error
	}
	return user, nil
}
