package models

import (
	"time"
)

type User struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Name      string    `gorm:"type:varchar(30);not null" json:"name"`
	Username  string    `gorm:"type:varchar(30);not null;unique;" josn:"username"`
	Email     string    `gorm:"type:varchar(50);not null; unique;" json:"email"`
	Password  string    `gorm:"type:varchar(100);not null;" json:"password"`
	Otp       int64     `gorm:"size:20;default:0" json:"otp"`
	Active    bool      `gorm:"default:false" json:"active_user"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
