package usecase

import "app/domain"

type BotRepository interface {
	FetchPublicOneById(string) (domain.PublicBot, error)
	FetchPublicAll() (domain.PublicBots, error)
	FetchOneById(string) (domain.Bot, error)
	FetchAll() (domain.Bots, error)
	Create(domain.Bot) (string, error)
}
