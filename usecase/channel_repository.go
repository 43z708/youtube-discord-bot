package usecase

import "app/domain"

type ChannelRepository interface {
	FetchOneById(string) (domain.Channel, error)
	FetchAllByBotID(string) (domain.Channels, error)
	FetchAllByGuildID(string) (domain.Channels, error)
	Create(domain.Channel) error
	Update(*domain.Channel) error
	Delete(string) error
}
