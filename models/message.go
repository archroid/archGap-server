package models

import (
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	SenderID      uint   `gorm:"not null"`
	ReceiverID    uint   `gorm:"not null"`
	Content       string `gorm:"not null"`
	Timestamp     string `gorm:"default:CURRENT_TIMESTAMP"`
	Status        string `gorm:"default:'sent'"`
	MessageType   string `gorm:"not null"`
	AttachmentURL string
}
