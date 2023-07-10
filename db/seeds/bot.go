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

	bot01 := domain.Bot{
		ID:           botId01,
		Name:         "test-bot01",
		TimeInterval: 3,
		Token:        botToken01,
		IsAvailable:  true,
	}
	err = db.Create(&bot01).Error

	botId02 := os.Getenv("BOT_ID02")
	botToken02 := os.Getenv("BOT_TOKEN02")

	bot02 := domain.Bot{
		ID:           botId02,
		Name:         "test-bot02",
		TimeInterval: 3,
		Token:        botToken02,
		IsAvailable:  true,
	}
	err = db.Create(&bot02).Error
	if err != nil {
		log.Fatalf("Error seeding bot data: %s", err.Error())
	}
	return nil
}
