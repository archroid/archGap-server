package config

import (
	"archroid/archGap/models"
	"archroid/archGap/utils"
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {

	utils.LoadEnv()

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	// Migrate the schema, create tables if they don't exist
	err = DB.AutoMigrate(&models.User{} , &models.Message{})
	if err != nil {
		panic("failed to migrate database: " + err.Error())
	}
}

func DbManager() *gorm.DB {
	return DB
}
