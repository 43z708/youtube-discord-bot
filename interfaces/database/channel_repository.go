package database

import (
	"app/domain"
	"strconv"

	"gorm.io/gorm"
)

type ChannelRepository struct {
	SqlHandler *gorm.DB
}

func (repo *ChannelRepository) FetchOneById(i string) (domain.Channel, error) {
	var channel domain.Channel
	id, _ := strconv.Atoi(i)
	result := repo.SqlHandler.First(&channel, id)
	if result.Error != nil {
		panic(result.Error)
	}

	return channel, result.Error
}

func (repo *ChannelRepository) FetchAllByBotID(botID string) (domain.Channels, error) {
	channels := make([]domain.Channel, 0)
	result := repo.SqlHandler.Where("bot_id = ?", botID).Find(&channels)
	if result.Error != nil {
		panic(result.Error)
	}
	var Channels domain.Channels
	for _, channel := range channels {
		Channel := domain.Channel{
			ID:             channel.ID,
			Name:           channel.Name,
			GuildID:        channel.GuildID,
			BotID:          channel.BotID,
			Searchword:     channel.Searchword,
			LastSearchedAt: channel.LastSearchedAt,
			CreatedAt:      channel.CreatedAt,
			UpdatedAt:      channel.UpdatedAt,
			DeletedAt:      channel.DeletedAt,
		}
		Channels = append(Channels, Channel)
	}

	return Channels, result.Error
}

func (repo *ChannelRepository) FetchAllByGuildID(guildID string) (domain.Channels, error) {
	channels := make([]domain.Channel, 0)
	result := repo.SqlHandler.Where("guild_id = ?", guildID).Find(&channels)
	if result.Error != nil {
		panic(result.Error)
	}
	var Channels domain.Channels
	for _, channel := range channels {
		Channel := domain.Channel{
			ID:             channel.ID,
			Name:           channel.Name,
			GuildID:        channel.GuildID,
			BotID:          channel.BotID,
			Searchword:     channel.Searchword,
			LastSearchedAt: channel.LastSearchedAt,
			CreatedAt:      channel.CreatedAt,
			UpdatedAt:      channel.UpdatedAt,
			DeletedAt:      channel.DeletedAt,
		}
		Channels = append(Channels, Channel)
	}

	return Channels, result.Error
}

func (repo *ChannelRepository) FetchAll() (domain.Channels, error) {
	channels := make([]domain.Channel, 0)
	result := repo.SqlHandler.Find(&channels)
	if result.Error != nil {
		panic(result.Error)
	}

	return channels, result.Error
}

func (repo *ChannelRepository) Create(g domain.Channel) error {

	result := repo.SqlHandler.Create(&g)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *ChannelRepository) Update(channel *domain.Channel) error {
	result := repo.SqlHandler.Save(&channel)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *ChannelRepository) Delete(i string) error {
	channel := make([]domain.Channel, 0)
	id, _ := strconv.Atoi(i)
	result := repo.SqlHandler.First(&channel, id)
	result = repo.SqlHandler.Delete(&channel)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
