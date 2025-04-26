package db

import (
	"errors"
	"slices"
	"otkritki/core/models"
)

var cardTypes = []string{
	"1", "2", "3", "4", "5", "6", "7", "8", "9", "10",
}

func (db *DB) GetCards() (*[]models.GiftCard, error) {
	var cards []models.GiftCard
	if err := db.db.Find(&cards).Error; err != nil {
		return nil, err
	}
	return &cards, nil
}

func (db *DB) AddCard(card *models.GiftCard) (cardId uint, err error) {
	if !slices.Contains[[]string, string](cardTypes, card.ImageType) {
		return 0, errors.New("Invalid card type")
	} else if err = db.db.Create(card).Error; err != nil {
		return 0, errors.New("Error adding card to database")
	}
	return card.ID, nil
}

func (db *DB) GetCardById(cardId uint) (*models.GiftCard, error) {
	var card models.GiftCard
	if err := db.db.First(&card, cardId).Error; err != nil {
		return nil, err
	}
	return &card, nil
}

func (db *DB) GetCardBySender(sender string) (*[]models.GiftCard, error) {
	var cards []models.GiftCard
	if err := db.db.Where(&models.GiftCard{From: sender}).Find(&cards).Error; err != nil {
		return nil, err
	}
	return &cards, nil
}

func (db *DB) GetCardByReceiver(recvName string) (*[]models.GiftCard, error) {
	var cards []models.GiftCard
	if err := db.db.Where(&models.GiftCard{To: recvName}).Find(&cards).Error; err != nil {
		return nil, err
	}
	return &cards, nil
}

func (db *DB) GetLatestCards(recvName string, number int) (*[]models.GiftCard, error) {
	var cards []models.GiftCard

	if err := db.db.Where(&models.GiftCard{To: recvName}).
		Order("date(created_at)").
		Limit(number).
		Find(&cards).Error; err != nil {

		return nil, err
	}
	return &cards, nil
}
