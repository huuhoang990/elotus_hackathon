package handlers

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	requests "go_api/cmd/api/request"
	"go_api/cmd/api/services"
	"go_api/cmd/common"
	"go_api/cmd/internal/models"
	"gorm.io/gorm"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func (h *Handler) RegisterHandler(c echo.Context) error {
	payload := new(requests.RegisterUserRequest)
	if err := (&echo.DefaultBinder{}).Bind(payload, c); err != nil {
		c.Logger().Error(err)
		return common.SendFailedValidationResponse(c, "bad request", err.Error())
	}

	c.Logger().Print(payload)
	validate := validator.New(validator.WithRequiredStructEnabled())
	validationErrors := validate.Struct(payload)
	if validationErrors != nil {
		c.Logger().Error(validationErrors)
		return common.SendFailedValidationResponse(c, "Failed", validationErrors.Error())
	}
	userService := services.NewUserService(h.DB)
	_, err := userService.GetUserByUserName(payload.UserName)
	if errors.Is(err, gorm.ErrRecordNotFound) == false {
		return common.SendFailedValidationResponse(c, "Failed", "User already exists")
	}
	registeredUser, err := userService.RegisterUser(payload)
	if err != nil {
		return common.SendFailedValidationResponse(c, "Failed", err.Error())
	}
	userInfo := map[string]interface{}{
		"id":        registeredUser.ID,
		"user_name": registeredUser.UserName,
	}
	return common.SendSuccessResponse(c, "User registered successfully", userInfo)
}

func (h *Handler) LoginHandler(c echo.Context) error {
	payload := new(requests.LoginRequest)
	if err := (&echo.DefaultBinder{}).BindBody(c, payload); err != nil {
		c.Logger().Error(err)
		return common.SendFailedValidationResponse(c, "bad request", err.Error())
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	validationErrors := validate.Struct(payload)
	if validationErrors != nil {
		c.Logger().Error(validationErrors)
		return common.SendFailedValidationResponse(c, "Failed", validationErrors.Error())
	}

	userService := services.NewUserService(h.DB)
	user, err := userService.GetUserByUserName(payload.UserName)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return common.SendFailedValidationResponse(c, "Failed", "Invalid username or password")
	}

	if common.ComparePasswordHash(payload.Password, user.Password) == false {
		return common.SendFailedValidationResponse(c, "Failed", "Invalid username or password")
	}

	accessToken, err := common.GenerateJWTToken(*user)
	return common.SendSuccessResponse(c, "Login successful", map[string]interface{}{
		"access_token": accessToken,
		"user_id":      user.ID,
	})
}

func (h *Handler) GetAuthenticatedUser(c echo.Context) error {
	user, ok := c.Get("user").(models.User)

	if !ok {
		return common.SendFailedValidationResponse(c, "Unauthorized", "User authentication failed")
	}
	return common.SendSuccessResponse(c, "Authenticated user", user)
}

func (h *Handler) UploadFile(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return common.SendFailedValidationResponse(c, "Failed", "Image upload error: "+err.Error())
	}

	if file.Size > 8*1024*1024 {
		return common.SendFailedValidationResponse(c, "Failed", "File size exceeds 8MB limit")
	}

	src, err := file.Open()
	if err != nil {
		return common.SendFailedValidationResponse(c, "Failed", "Image open error: "+err.Error())
	}
	defer src.Close()

	// Read first 512 bytes for content type detection
	buf := make([]byte, 512)
	n, err := src.Read(buf)
	if err != nil && err != io.EOF {
		return common.SendFailedValidationResponse(c, "Failed", "Image read error: "+err.Error())
	}

	contentType := http.DetectContentType(buf[:n])
	allowedTypes := []string{"image/jpeg", "image/png", "image/gif"}

	validType := false
	for _, t := range allowedTypes {
		if t == contentType {
			validType = true
			break
		}
	}

	if !validType {
		return common.SendFailedValidationResponse(c, "Failed", "Invalid file type. Only JPEG, PNG, and GIF are allowed.")
	}

	tempPath := filepath.Join("/tmp", file.Filename)
	dst, err := os.Create(tempPath)
	if err != nil {
		return common.SendFailedValidationResponse(c, "Failed", "Image create error: "+err.Error())
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return common.SendFailedValidationResponse(c, "Failed", "Image copy error: "+err.Error())
	}

	imageInfoService := services.NewImageInfoService(h.DB)
	user, ok := c.Get("user").(models.User)
	if !ok {
		return common.SendFailedValidationResponse(c, "Unauthorized", "User authentication failed")
	}
	imageInfo, err := imageInfoService.SaveImageInfo(user.ID, file.Filename, contentType, file.Size)
	if err != nil {
		return common.SendFailedValidationResponse(c, "Failed", "Image info save error: "+err.Error())
	}

	return common.SendSuccessResponse(c, "File uploaded successfully", imageInfo)
}
