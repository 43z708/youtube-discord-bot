package main

import (
	"app/domain"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func BotSeeds(db *gorm.DB) error {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading dotenv: %s", err.Error())
	}

	botId01 := os.Getenv("BOT_ID01")
	botToken01 := os.Getenv("BOT_TOKEN01")

	bot := domain.Bot{
		ID:          botId01,
		Name:        "test-bot01",
		Token:       botToken01,
		IsAvailable: false,
	}
	err = db.Create(&bot).Error
	if err != nil {
		log.Fatalf("Error seeding bot data: %s", err.Error())
	}
	return nil
}
