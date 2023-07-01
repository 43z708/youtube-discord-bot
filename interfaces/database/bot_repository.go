package database

import (
	"app/domain"
	"strconv"

	"gorm.io/gorm"
)

type BotRepository struct {
	SqlHandler *gorm.DB
}

func (repo *BotRepository) FetchPublicOneById(i string) (domain.PublicBot, error) {
	var bot domain.Bot
	id, _ := strconv.Atoi(i)
	result := repo.SqlHandler.First(&bot, id)
	if result.Error != nil {
		panic(result.Error)
	}
	publicBot := domain.PublicBot{
		ID:           bot.ID,
		Name:         bot.Name,
		TimeInterval: bot.TimeInterval,
		IsAvailable:  bot.IsAvailable,
	}

	return publicBot, result.Error
}

func (repo *BotRepository) FetchPublicAll() (domain.PublicBots, error) {
	bots := make([]domain.Bot, 0)
	result := repo.SqlHandler.Find(&bots)
	if result.Error != nil {
		panic(result.Error)
	}
	var publicBots domain.PublicBots
	for _, bot := range bots {
		publicBot := domain.PublicBot{
			ID:           bot.ID,
			Name:         bot.Name,
			TimeInterval: bot.TimeInterval,
			IsAvailable:  bot.IsAvailable,
		}
		publicBots = append(publicBots, publicBot)
	}

	return publicBots, result.Error
}

func (repo *BotRepository) FetchOneById(i string) (domain.Bot, error) {
	var bot domain.Bot
	id, _ := strconv.Atoi(i)
	result := repo.SqlHandler.First(&bot, id)
	if result.Error != nil {
		panic(result.Error)
	}
	return bot, result.Error
}

func (repo *BotRepository) FetchAll() (domain.Bots, error) {
	bots := make([]domain.Bot, 0)
	result := repo.SqlHandler.Find(&bots)
	if result.Error != nil {
		panic(result.Error)
	}

	return bots, result.Error
}

func (repo *BotRepository) Create(b domain.Bot) error {

	result := repo.SqlHandler.Create(&b)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *BotRepository) Update(b *domain.Bot) error {
	result := repo.SqlHandler.Save(&b)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// func (repo *BotRepository) DeleteBot(id int) error {
// 	bot := make([]domain.Bot, 0)
// 	result := repo.SqlHandler.Conn.First(&bot, id)
// 	if result.Error != nil {
// 		panic(result.Error)
// 	}
// 	repo.SqlHandler.Conn.Delete(&bot)
// 	return nil
// }
