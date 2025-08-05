package main

import (
	"go_api/cmd/common"
	"go_api/cmd/internal/models"
	"log"
)

func main() {
	db, err := common.NewMysql()
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&models.ImageInfo{})
	if err != nil {
		panic(err)
	}
	log.Println("Database migrated")
}
