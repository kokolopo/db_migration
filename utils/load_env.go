package utils

import (
	"log"

	"github.com/joho/godotenv"
)

// Fungsi untuk memuat environment dari file .env
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
