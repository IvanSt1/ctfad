package db

import (
	"errors"
	"otkritki/core/models"
)

func (db *DB) AddUser(user *models.User) (id uint, err error) {
	if _, err := db.GetUserByName(user.Username); err == nil {
		return 0, errors.New("User alerady exists")
	}
	if err := db.db.Create(user).Error; err != nil {
		return 0, errors.New("Could not add user")
	}
	return user.ID, nil
}

func (db *DB) GetUserById(id uint) (*models.User, error) {
	var result models.User
	if err := db.db.First(&result, id).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

func (db *DB) GetUserByName(name string) (*models.User, error) {
	var result models.User
	if err := db.db.Where(
		&models.User{
			Username: name,
		}).First(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

func (db *DB) GetUserByNameAndPassword(name, pass string) (*models.User, error) {
	var result models.User
	if err := db.db.Where(
		&models.User{
			Username: name,
			Password: pass,
		}).Find(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}
