package services

import (
	"archroid/archGap/models"
	"archroid/archGap/utils"
	"errors"

	"gorm.io/gorm"
)

// LoginUser checks if the user exists, compares the password, and returns the user if valid
func LoginUser(db *gorm.DB, email, password string) (*models.User, string, error) {
	var user models.User

	// Find the user by email
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
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

// RegisterUser handles the registration of a new user
func RegisterUser(db *gorm.DB, email, password string) (*models.User, error) {
	// Check if user with the same email already exists
	var existingUser models.User
	if err := db.Where("email = ?", email).First(&existingUser).Error; err == nil {
		// If user already exists, return an error
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
	if err := db.Create(&newUser).Error; err != nil {
		return nil, err
	}

	// Return the created user (excluding the password)
	return &newUser, nil
}

// UpdateProfile updates the user's email and/or password
func UpdateProfile(db *gorm.DB, userID uint, username string, profilepicture string) (*models.User, error) {
	var user models.User

	// Find the user by their ID
	if err := db.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	if username != "" {
		user.Name = username
	}

	if profilepicture != "" {
		user.ProfilePicture = profilepicture
	}

	if err := db.Save(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
