package services

import (
	"archroid/archGap/models"

	"gorm.io/gorm"
)

// SaveMessage saves a message with or without an attachment to the database
func SaveMessage(db *gorm.DB, message models.Message) error {
	// Start a transaction
	tx := db.Begin()

	// Save the message
	if err := tx.Create(&message).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	tx.Commit()
	return nil
}
