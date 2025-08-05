package models

import "gorm.io/gorm"

type ImageInfo struct {
	gorm.Model
	UserId      uint   `gorm:"not null" json:"user_id"`
	FileName    string `gorm:"not null" json:"file_name"`
	ContentType string `gorm:"not null" json:"content_type"`
	Size        uint   `gorm:"not null" json:"size"`
}
