package controllers

import (
	"app/domain"
	"app/interfaces/database"
	"app/usecase"
	"log"
	"strconv"

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

func (controller *BotController) Index() domain.Bots {
	bots, err := controller.Interactor.Bots()
	if err != nil {
		log.Fatalf("Error getting bots data: %s", err.Error())
	}

	return bots
}

func (controller *BotController) Show(i string) domain.Bot {
	id, _ := strconv.Atoi(i)
	bot, err := controller.Interactor.BotById(id)

	if err != nil {
		log.Fatalf("Error getting bots data: %s", err.Error())
	}
	return bot
}
