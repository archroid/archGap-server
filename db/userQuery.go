package db

import (
	"archroid/archGap/models"
	"archroid/archGap/utils"
	"errors"
	"time"

	"gorm.io/gorm"
)

func LoginUser(email string, password string) (*models.User, string, error) {
	var user models.User

	// Find the user by email
	if err := DB.Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, "", errors.New("user not found")
		}
		return nil, "", err
	}

	// Compare the provided password with the stored hashed password
	if !utils.ComparePassword(password, user.Password) {
		return nil, "", errors.New("invalid password")
	}

	// Generate JWT token
	token, err := utils.CreateJWT(user.ID)
	if err != nil {
		return nil, "", err
	}

	return &user, token, nil
}

func RegisterUser(email string, password string) (*models.User, error) {
	// Check if user with the same email already exists
	var existingUser models.User
	if err := DB.Where("email = ?", email).First(&existingUser).Error; err == nil {
		return nil, errors.New("email is already in use")
	}
	// Hash the password before saving it to the database
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	// Create a new user
	newUser := models.User{
		Email:    email,
		Password: hashedPassword,
	}

	// Save the user in the database
	if err := DB.Create(&newUser).Error; err != nil {
		return nil, err
	}

	// Return the created user
	return &newUser, nil
}

func GetUser(userID uint) (*models.User, error) {
	user := models.User{}
	err := DB.Where("id = ?", userID).First(&user).Error
	if err != nil {
		return nil, errors.New("error getting user")
	}
	return &user, nil
}

func UpdateUser(userID uint, name *string, profilePicture *string, lastseen *time.Time, isonline *bool) error {
	updates := map[string]interface{}{}

	if name != nil {
		updates["name"] = *name
	}
	if profilePicture != nil {
		updates["profile_picture"] = *profilePicture
	}
	if isonline != nil {
		updates["is_online"] = *isonline
	}
	if lastseen != nil {
		updates["last_seen"] = *lastseen
	}

	err := DB.Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error
	if err != nil {
		return errors.New("error updating user")
	}
	return nil
}


func GetChatsbyUser(userID uint) ([]models.Chat, error) {
	var chats []models.Chat

	err := DB.Joins("JOIN chat_participants ON chat_participants.chat_id = chats.id").
		Where("chat_participants.user_id = ?", userID).
		Find(&chats).Error

	if err != nil {
		return nil, errors.New("error getting chats")
	}

	return chats, nil
}


// return if user is online or not
func GetUserOnline(userID uint) (bool, error) {
	var isOnline bool
	err := DB.Model(&models.User{}).Where("id = ?", userID).Select("is_online").Scan(&isOnline).Error
	if err != nil {
		return false, errors.New("error getting user")
	}
	return isOnline, nil
}


