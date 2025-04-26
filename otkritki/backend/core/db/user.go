package db

import (
    "github.com/IvanSt1/ctfad/otkritki/backend/core/models"
)

// AddUser создаёт нового пользователя в базе и возвращает его с заполненным ID.
func AddUser(user *models.User) (*models.User, error) {
    if err := GetDB().Create(user).Error; err != nil {
        return nil, err
    }
    return user, nil
}

// GetUserById ищет пользователя по ID.
func GetUserById(id uint) (*models.User, error) {
    var u models.User
    if err := GetDB().First(&u, id).Error; err != nil {
        return nil, err
    }
    return &u, nil
}

// GetUserByName ищет пользователя по имени.
func GetUserByName(username string) (*models.User, error) {
    var u models.User
    if err := GetDB().Where("username = ?", username).First(&u).Error; err != nil {
        return nil, err
    }
    return &u, nil
}
