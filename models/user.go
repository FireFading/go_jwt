package models

type User struct {
	Id       uint   `json:"id" gorm:"primary_key"`
	Name     string `json:"name"`
	Email    string `gorm:"unique"`
	Password []byte `json:"-"`
}
