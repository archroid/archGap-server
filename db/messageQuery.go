package db

import (
	"archroid/archGap/models"
	"errors"
)

func SendMessage(chatID uint, senderID uint, content string, messageType string) error {
	message := models.Message{
		SenderID:    senderID,
		ChatID:      chatID,
		Content:     content,
		MessageType: messageType,
		Status:      "sent",
	}

	err := DB.Create(&message).Error
	if err != nil {
		return errors.New("error sending message")
	}
	return nil
}

// return all messages of a chat
func GetMessagesinChat(chatID uint) ([]models.Message, error) {
	messages := []models.Message{}
	err := DB.Where("chat_id = ?", chatID).Find(&messages).Error
	if err != nil {
		return nil, errors.New("error getting messages")
	}
	return messages, nil
}
