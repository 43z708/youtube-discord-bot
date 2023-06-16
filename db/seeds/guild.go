package main

import (
	"app/domain"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func GuildSeeds(db *gorm.DB) error {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading dotenv: %s", err.Error())
	}

	guildId01 := os.Getenv("GUILD_ID01")
	botId01 := os.Getenv("BOT_ID01")

	guild := domain.Guild{
		ID:    guildId01,
		Name:  "test-guild01",
		BotID: botId01,
	}
	err = db.Create(&guild).Error
	if err != nil {
		log.Fatalf("Error seeding guild data: %s", err.Error())
	}
	return nil
}
