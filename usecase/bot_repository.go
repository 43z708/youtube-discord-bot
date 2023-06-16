package usecase

import "app/domain"

type BotRepository interface {
	Store(domain.Bot) (string, error)
	FindPublicById(int) (domain.PublicBot, error)
	FindPublicAll() (domain.PublicBots, error)
	FindById(int) (domain.Bot, error)
	FindAll() (domain.Bots, error)
}
