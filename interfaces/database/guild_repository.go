package database

import (
	"app/domain"
	"strconv"

	"gorm.io/gorm"
)

type GuildRepository struct {
	SqlHandler *gorm.DB
}

func (repo *GuildRepository) FetchPublicOneById(i string) (domain.PublicGuild, error) {
	var guild domain.Guild
	id, _ := strconv.Atoi(i)
	result := repo.SqlHandler.First(&guild, id)
	if result.Error != nil {
		panic(result.Error)
	}
	publicGuild := domain.PublicGuild{
		ID:    guild.ID,
		Name:  guild.Name,
		BotID: guild.BotID,
	}

	return publicGuild, result.Error
}

func (repo *GuildRepository) FetchPublicAllByBotID(botID string) (domain.PublicGuilds, error) {
	guilds := make([]domain.Guild, 0)
	result := repo.SqlHandler.Where("bot_id = ?", botID).Find(&guilds)
	if result.Error != nil {
		panic(result.Error)
	}
	var publicGuilds domain.PublicGuilds
	for _, guild := range guilds {
		publicGuild := domain.PublicGuild{
			ID:    guild.ID,
			Name:  guild.Name,
			BotID: guild.BotID,
		}
		publicGuilds = append(publicGuilds, publicGuild)
	}

	return publicGuilds, result.Error
}

func (repo *GuildRepository) FetchOneById(i string) (domain.Guild, error) {
	var guild domain.Guild
	id, _ := strconv.Atoi(i)
	result := repo.SqlHandler.First(&guild, id)
	if result.Error != nil {
		panic(result.Error)
	}
	return guild, result.Error
}

func (repo *GuildRepository) FetchAll() (domain.Guilds, error) {
	guilds := make([]domain.Guild, 0)
	result := repo.SqlHandler.Find(&guilds)
	if result.Error != nil {
		panic(result.Error)
	}

	return guilds, result.Error
}

func (repo *GuildRepository) Create(g domain.Guild) (string, error) {

	result := repo.SqlHandler.Create(&g)
	if result.Error != nil {
		return "", result.Error
	}
	return g.ID, result.Error
}

func (repo *GuildRepository) Update(guild *domain.Guild) (string, error) {
	result := repo.SqlHandler.Save(&guild)
	if result.Error != nil {
		return "", result.Error
	}
	return guild.ID, result.Error
}

func (repo *GuildRepository) Delete(i string) (string, error) {
	guild := make([]domain.Guild, 0)
	id, _ := strconv.Atoi(i)
	result := repo.SqlHandler.First(&guild, id)
	if result.Error != nil {
		return "", result.Error
	}
	result = repo.SqlHandler.Delete(&guild)
	return i, result.Error
}
