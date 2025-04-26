package db

import (
    "github.com/IvanSt1/ctfad/otkritki/backend/core/models"
)

// ListCards возвращает все открытки из БД
func ListCards() ([]models.GiftCard, error) {
    var cards []models.GiftCard
    if err := GetDB().Find(&cards).Error; err != nil {
        return nil, err
    }
    return cards, nil
}

// CreateCard сохраняет новую открытку в БД
func CreateCard(card *models.GiftCard) error {
    return GetDB().Create(card).Error
}
