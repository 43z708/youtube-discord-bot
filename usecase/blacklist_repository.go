package usecase

import "app/domain"

type BlacklistRepository interface {
	FetchOneById(string) (domain.Blacklist, error)
	FetchOneByDistributor(string) (domain.Blacklist, error)
	FetchAllByGuildID(string) (domain.Blacklists, error)
	Create(domain.Blacklist) error
	Update(*domain.Blacklist) error
	Delete(int) error
}
