package database

import (
	"github.com/inadislam/go-login-api/models"
	"github.com/inadislam/go-login-api/utils"
)

// var pass, err = utils.HashPassword("1234")
// var users = models.User{
// 	Name:     "Nazim Uddin",
// 	Username: "inadislam",
// 	Email:    "inadislam@gmail.com",
// 	Password: string(pass),
// }

func AutoMigrator() {
	db := Connect()
	defer db.Close()

	// err := db.Debug().DropTableIfExists(&models.User{}).Error
	// utils.CheckErr(err)

	err := db.Debug().AutoMigrate(&models.User{}).Error
	utils.CheckErr(err)

	// err = db.Debug().Create(&users).Error
	// utils.CheckErr(err)
}
