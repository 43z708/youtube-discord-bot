package usecase

import "app/domain"

type BlacklistInteractor struct {
	BlacklistRepository BlacklistRepository
}

func (interactor *BlacklistInteractor) FetchsById(identifier string) (blacklist domain.Blacklist, err error) {
	blacklist, err = interactor.BlacklistRepository.FetchOneById(identifier)
	return blacklist, err
}

func (interactor *BlacklistInteractor) FetchAllByBotID(identifier string) (channels domain.Blacklists, err error) {
	channels, err = interactor.BlacklistRepository.FetchAllByBotID(identifier)
	return channels, err
}

func (interactor *BlacklistInteractor) FetchAllByGuildID(identifier string) (channels domain.Blacklists, err error) {
	channels, err = interactor.BlacklistRepository.FetchAllByGuildID(identifier)
	return channels, err
}

func (interactor *BlacklistInteractor) Create(b domain.Blacklist) (err error) {
	_, err = interactor.BlacklistRepository.Create(b)
	return err
}

func (interactor *BlacklistInteractor) Update(b domain.Blacklist) (err error) {
	_, err = interactor.BlacklistRepository.Update(b)
	return err
}
func (interactor *BlacklistInteractor) Delete(identifier string) (err error) {
	_, err = interactor.BlacklistRepository.Delete(identifier)
	return err
}
