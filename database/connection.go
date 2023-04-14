package database

import (
	"go_jwt/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	connection, err := gorm.Open(mysql.Open("appuser:secretpassword@tcp(mysql:3306)/app"), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	DB = connection

	connection.AutoMigrate(&models.User{})
}
