package services

import (
	"go_api/cmd/internal/models"
	"gorm.io/gorm"
)

type ImageInfoService struct {
	db *gorm.DB
}

func NewImageInfoService(db *gorm.DB) *ImageInfoService {
	return &ImageInfoService{
		db: db,
	}
}

func (imageInfoService *ImageInfoService) SaveImageInfo(userId uint, fileName string, contentType string, fileSize int64) (*models.ImageInfo, error) {
	imageInfo := models.ImageInfo{
		UserId:      userId,
		FileName:    fileName,
		ContentType: contentType,
		Size:        uint(fileSize),
	}
	result := imageInfoService.db.Create(&imageInfo)
	if result.Error != nil {
		return nil, result.Error
	}
	return &imageInfo, nil
}
