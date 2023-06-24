package controllers

import (
	"app/domain"
	"app/interfaces/database"
	"app/usecase"
	"log"

	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

type GuildController struct {
	GuildInteractor usecase.GuildInteractor
}

func NewGuildController(sqlHandler *gorm.DB) *GuildController {
	return &GuildController{
		GuildInteractor: usecase.GuildInteractor{
			GuildRepository: &database.GuildRepository{
				SqlHandler: sqlHandler,
			},
		},
	}
}

func (controller *GuildController) FetchPublicOneById(id string) domain.PublicGuild {
	guild, err := controller.GuildInteractor.FetchPublicOneById(id)

	if err != nil {
		log.Fatalf("Error getting guilds data: %s", err.Error())
	}
	return guild
}

func (controller *GuildController) FetchPublicAllByBotID(id string) domain.PublicGuilds {
	guilds, err := controller.GuildInteractor.FetchPublicAllByBotID(id)
	if err != nil {
		log.Fatalf("Error getting channels data: %s", err.Error())
	}

	return guilds
}

func (controller *GuildController) FetchOneById(id string) domain.Guild {
	guild, err := controller.GuildInteractor.FetchOneById(id)

	if err != nil {
		log.Fatalf("Error getting guilds data: %s", err.Error())
	}
	return guild
}

func (controller *GuildController) FetchAll() domain.Guilds {
	guilds, err := controller.GuildInteractor.FetchAll()
	if err != nil {
		log.Fatalf("Error getting guilds data: %s", err.Error())
	}

	return guilds
}

func (controller *GuildController) Create(s *discordgo.Session, i *discordgo.InteractionCreate) {

}

func (controller *GuildController) Update(s *discordgo.Session, c *discordgo.GuildUpdate) {

}

func (controller *GuildController) Delete(s *discordgo.Session, c *discordgo.ChannelDelete) {
}
