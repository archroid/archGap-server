package db

import (
	"archroid/archGap/models"
	"errors"
)

func NewChat(chatname string) (*models.Chat, error) {
	chat := models.Chat{ChatName: chatname}
	err := DB.Create(&chat).Error
	if err != nil {
		return nil, errors.New("error creating chat")
	}
	return &chat, nil
}

func AddUserToChat(chatID uint, userIDs []uint) ([]models.ChatParticipant, error) {
	participants := []models.ChatParticipant{}
	for _, userID := range userIDs {
		participant := models.ChatParticipant{ChatID: chatID, UserID: userID}
		participants = append(participants, participant)
	}
	err := DB.Create(&participants).Error
	if err != nil {
		return nil, errors.New("error adding user to chat")
	}
	return participants, nil
}

func IsChatExist(userID1 uint, userID2 uint) (bool, uint, error) {
	var chatParticipant models.ChatParticipant
	var chat models.Chat
	err := DB.
		Where("is_group = ?", false).
		Joins("JOIN chat_participants cp1 ON cp1.chat_id = chats.id").
		Joins("JOIN chat_participants cp2 ON cp2.chat_id = chats.id").
		Where("cp1.user_id = ? AND cp2.user_id = ?", userID1, userID2).
		First(&chat).Error

	if err != nil {
		return false, 0, errors.New("error finding chat")
	}
	return true, chatParticipant.ChatID, nil
}


// return all users of a chat
func GetUsersbyChat(chatID uint) ([]models.User, error) {
	users := []models.User{}
	err := DB.Joins("JOIN chat_participants ON chat_participants.user_id = users.user_id").
		Where("chat_participants.chat_id = ?", chatID).
		Find(&users).Error
	if err != nil {
		return nil, errors.New("error getting users")
	}
	return users, nil
}
