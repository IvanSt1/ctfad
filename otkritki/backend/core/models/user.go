package models

import (
    "errors"

    "golang.org/x/crypto/bcrypt"    // ИЗМЕНЕНО
    "gorm.io/gorm"
)

type User struct {
    ID           uint   `gorm:"primaryKey"`
    Username     string `gorm:"uniqueIndex;not null"`
    PasswordHash string `gorm:"not null"`     // ИЗМЕНЕНО: переименовано из Password
}

// Создаёт нового пользователя с хешированным паролем
func CreateUser(username, password string) error {
    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    user := User{Username: username, PasswordHash: string(hash)}
    return db.Create(&user).Error
}

// Находит пользователя по имени
func FindUserByUsername(username string) (*User, error) {
    var u User
    if err := db.Where("username = ?", username).First(&u).Error; err != nil {
        return nil, err
    }
    return &u, nil
}

// Сравнивает хеш с введённым паролем
func ComparePassword(hash, password string) error {
    if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
        return errors.New("password mismatch")
    }
    return nil
}
