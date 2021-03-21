package database

import (
	"github.com/inadislam/go-login-api/utils"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func Connect() *gorm.DB {
	url := "host=localhost port=5432 user=inadislam password=root dbname=api"
	db, err := gorm.Open("postgres", url)
	utils.CheckErr(err)
	return db
}
