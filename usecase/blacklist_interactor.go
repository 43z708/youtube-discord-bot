package usecase

import "app/domain"

type BlacklistInteractor struct {
	BlacklistRepository BlacklistRepository
}

func (interactor *BlacklistInteractor) FetchOneById(identifier string) (blacklist domain.Blacklist, err error) {
	blacklist, err = interactor.BlacklistRepository.FetchOneById(identifier)
	return blacklist, err
}
func (interactor *BlacklistInteractor) FetchOneByDistributor(distributor string) (blacklist domain.Blacklist, err error) {
	blacklist, err = interactor.BlacklistRepository.FetchOneByDistributor(distributor)
	return blacklist, err
}

func (interactor *BlacklistInteractor) FetchAllByGuildID(identifier string) (blacklists domain.Blacklists, err error) {
	blacklists, err = interactor.BlacklistRepository.FetchAllByGuildID(identifier)
	return blacklists, err
}

func (interactor *BlacklistInteractor) Create(b domain.Blacklist) (err error) {
	err = interactor.BlacklistRepository.Create(b)
	return err
}

func (interactor *BlacklistInteractor) Update(b *domain.Blacklist) (err error) {
	err = interactor.BlacklistRepository.Update(b)
	return err
}
func (interactor *BlacklistInteractor) Delete(identifier int) (err error) {
	err = interactor.BlacklistRepository.Delete(identifier)
	return err
}
