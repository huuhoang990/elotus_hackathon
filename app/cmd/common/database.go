package common

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewMysql() (*gorm.DB, error) {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
		return nil, err
	}

	host := os.Getenv("MYSQL_DB_CONNECTION")
	port := os.Getenv("MYSQL_PORT")
	database := os.Getenv("MYSQL_DATABASE")
	userName := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", userName, password, host, port, database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	log.Default().Println("Successfully connected to database")
	return db, nil
}
