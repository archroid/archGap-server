package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email          string `gorm:"uniqueIndex;not null"`
	Password       string
	Name           string
	ProfilePicture string
	LastSeen       *time.Time `gorm:"index"`
	IsOnline       bool

	Chats    []Chat    `gorm:"many2many:chat_participants"`
	Messages []Message `gorm:"foreignKey:SenderID"` // ✅ Explicit foreign key
}

type Message struct {
	gorm.Model
	SenderID    uint   `gorm:"not null"` // The sender of the message
	ChatID      uint   `gorm:"not null"` // The chat where the message was sent
	Content     string `gorm:"type:text"`
	Status      string `gorm:"default:'sent'"`
	MessageType string // e.g., "text", "image", "video"

	Sender User `gorm:"foreignKey:SenderID;constraint:OnDelete:CASCADE"`
	Chat   Chat `gorm:"foreignKey:ChatID;constraint:OnDelete:CASCADE"`
}

type Chat struct {
	gorm.Model
	ChatName string
	IsGroup  bool
	Avatar   string

	Participants []User `gorm:"many2many:chat_participants"`
	Messages     []Message
}

type ChatParticipant struct {
	ChatID uint `gorm:"primaryKey"`
	UserID uint `gorm:"primaryKey"`

	// Foreign keys
	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Chat Chat `gorm:"foreignKey:ChatID;constraint:OnDelete:CASCADE"`
}
