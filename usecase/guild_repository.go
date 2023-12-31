package usecase

import "app/domain"

type GuildRepository interface {
	FetchPublicOneById(string) (domain.PublicGuild, error)
	FetchPublicAllByBotID(string) (domain.PublicGuilds, error)
	FetchOneById(string) (domain.Guild, error)
	FetchAll() (domain.Guilds, error)
	Create(domain.Guild) error
	Update(*domain.Guild) error
	Delete(string) error
}
