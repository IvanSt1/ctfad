package models

import "golang.org/x/crypto/bcrypt"

// Gender — тип пола пользователя.
type Gender string

const (
    Male   Gender = "male"
    Female Gender = "female"
)

// User описывает модель пользователя.
type User struct {
    ID           uint   `gorm:"primaryKey"`
    Username     string `gorm:"uniqueIndex;not null"`
    PasswordHash string `gorm:"not null"`
    Gender       Gender `gorm:"type:ENUM('male','female');not null"`
}

// ComparePassword сравнивает хеш и пароль.
func ComparePassword(hash, password string) error {
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
