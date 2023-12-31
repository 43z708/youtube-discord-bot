package controllers

import (
	"app/domain"
	"app/interfaces/database"
	"app/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BotApiController struct {
	BotInteractor usecase.BotInteractor
}

func NewBotApiController(sqlHandler *gorm.DB) *BotApiController {
	return &BotApiController{
		BotInteractor: usecase.BotInteractor{
			BotRepository: &database.BotRepository{
				SqlHandler: sqlHandler,
			},
		},
	}
}

func (controller *BotApiController) FetchPublicOneById(c Context) {
	id := c.Param("id")
	bot, err := controller.BotInteractor.FetchPublicOneById(id)
	if err != nil {
		c.JSON(500, err.Error())
	}
	c.JSON(200, bot)
}

func (controller *BotApiController) FetchAllPublic(c Context) {
	bots, err := controller.BotInteractor.FetchAllPublic()
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, bots)
}

func (controller *BotApiController) Create(c Context) {
	b := domain.Bot{}
	c.Bind(&b)
	err := controller.BotInteractor.Create(b)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	c.JSON(201, "ok")
}
