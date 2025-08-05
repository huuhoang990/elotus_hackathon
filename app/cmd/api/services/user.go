package services

import (
	requests "go_api/cmd/api/request"
	"go_api/cmd/common"
	"go_api/cmd/internal/models"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		db: db,
	}
}

func (userService UserService) RegisterUser(userRequest *requests.RegisterUserRequest) (*models.User, error) {
	hashPassword, err := common.HashPassword(userRequest.Password)
	if err != nil {
		return nil, err
	}
	createdUser := models.User{
		UserName: userRequest.UserName,
		Password: hashPassword,
	}
	result := userService.db.Create(&createdUser)
	if result.Error != nil {
		return nil, result.Error
	}

	return &createdUser, nil
}

func (userService UserService) GetUserByUserName(userName string) (*models.User, error) {
	var user models.User
	if err := userService.db.Where("user_name = ?", userName).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
