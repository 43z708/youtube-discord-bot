package controllers

import (
	"app/domain"
	"app/interfaces/database"
	"app/usecase"
	"log"

	"gorm.io/gorm"
)

type BotController struct {
	Interactor usecase.BotInteractor
}

func NewBotController(sqlHandler *gorm.DB) *BotController {
	return &BotController{
		Interactor: usecase.BotInteractor{
			BotRepository: &database.BotRepository{
				SqlHandler: sqlHandler,
			},
		},
	}
}

func (controller *BotController) FetchOneById(id string) domain.Bot {
	bot, err := controller.Interactor.FetchOneById(id)

	if err != nil {
		log.Fatalf("Error getting bots data: %s", err.Error())
	}
	return bot
}

func (controller *BotController) FetchAll() domain.Bots {
	bots, err := controller.Interactor.FetchAll()
	if err != nil {
		log.Fatalf("Error getting bots data: %s", err.Error())
	}

	return bots
}
