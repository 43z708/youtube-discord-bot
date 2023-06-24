package usecase

import "app/domain"

type BlacklistRepository interface {
	FetchOneById(string) (domain.Blacklist, error)
	FetchAllByBotID(string) (domain.Blacklists, error)
	FetchAllByGuildID(string) (domain.Blacklists, error)
	Create(domain.Blacklist) (string, error)
	Update(domain.Blacklist) (domain.Blacklists, error)
	Delete(string) (domain.Blacklists, error)
}
