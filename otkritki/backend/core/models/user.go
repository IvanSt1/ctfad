package models 

import "gorm.io/gorm"

type Gender string

const (
	Male   Gender = "male"
	Female Gender = "female"
)

type User struct {
	gorm.Model `json:"-"`
	Username   string `json:"username"`
	Password   string `json:"-"`
	Gender     Gender `json:"gender"`
}

func NewUser(username, password string, gender Gender) *User {
	return &User{
		Username: username,
		Password: password,
		Gender:   gender,
	}
}
