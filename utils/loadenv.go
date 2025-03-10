package utils

import (
	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Error("Error loading .env file")
	}
}