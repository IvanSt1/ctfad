package db

import (
    "gorm.io/gorm"

    // Правильный абсолютный путь модуля
    "github.com/IvanSt1/ctfad/otkritki/backend/core/models"
)

type Card struct {
    gorm.Model
    To   string
    Text string
}

// ListCards возвращает все записи из таблицы cards
func ListCards(db *gorm.DB) ([]models.Card, error) {
    var cards []models.Card
    if err := db.Find(&cards).Error; err != nil {
        return nil, err
    }
    return cards, nil
}

// CreateCard сохраняет новую карточку
func CreateCard(db *gorm.DB, card *models.Card) error {
    return db.Create(card).Error
}
