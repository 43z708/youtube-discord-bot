package database

import (
	"app/domain"

	"gorm.io/gorm"
)

type BotRepository struct {
	SqlHandler *gorm.DB
}

func (repo *BotRepository) Store(b domain.Bot) (string, error) {

	result := repo.SqlHandler.Create(&b)
	if result.Error != nil {
		return "", result.Error
	}
	return b.ID, result.Error
}

func (repo *BotRepository) FindPublicAll() (domain.PublicBots, error) {
	bots := make([]domain.Bot, 0)
	result := repo.SqlHandler.Find(&bots)
	if result.Error != nil {
		panic(result.Error)
	}
	var publicBots domain.PublicBots
	for _, bot := range bots {
		publicBot := domain.PublicBot{
			ID:          bot.ID,
			Name:        bot.Name,
			IsAvailable: bot.IsAvailable,
		}
		publicBots = append(publicBots, publicBot)
	}

	return publicBots, result.Error
}

func (repo *BotRepository) FindPublicById(id int) (domain.PublicBot, error) {
	var bot domain.Bot
	result := repo.SqlHandler.First(&bot, id)
	if result.Error != nil {
		panic(result.Error)
	}
	publicBot := domain.PublicBot{
		ID:          bot.ID,
		Name:        bot.Name,
		IsAvailable: bot.IsAvailable,
	}

	return publicBot, result.Error
}

func (repo *BotRepository) FindAll() (domain.Bots, error) {
	bots := make([]domain.Bot, 0)
	result := repo.SqlHandler.Find(&bots)
	if result.Error != nil {
		panic(result.Error)
	}

	return bots, result.Error
}

func (repo *BotRepository) FindById(id int) (domain.Bot, error) {
	var bot domain.Bot
	result := repo.SqlHandler.First(&bot, id)
	if result.Error != nil {
		panic(result.Error)
	}
	return bot, result.Error
}

// func (repo *BotRepository) UpdateBot(bot *domain.Bot) error {
// 	repo.SqlHandler.Conn.Save(&bot)
// 	return nil
// }

// func (repo *BotRepository) DeleteBot(id int) error {
// 	bot := make([]domain.Bot, 0)
// 	result := repo.SqlHandler.Conn.First(&bot, id)
// 	if result.Error != nil {
// 		panic(result.Error)
// 	}
// 	repo.SqlHandler.Conn.Delete(&bot)
// 	return nil
// }
