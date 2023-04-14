package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Connect() {
	_, err := gorm.Open(mysql.Open("appuser:secretpassword@tcp(mysql:3306)/app"), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}
}
