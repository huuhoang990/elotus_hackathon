package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserName string `gorm:"type:varchar(200);not null;uniqueIndex" json:"user_name"`
	Password string `gorm:"type:varchar(200);not null" json:"password"`
}
