package models 

import "gorm.io/gorm"

type GiftCard struct {
	*gorm.Model `json:"-"`
	To          string `json:"to"`
	From        string `json:"from"`
	Text        string `json:"text"`
	ImageType   string `json:"imageType"`
}
