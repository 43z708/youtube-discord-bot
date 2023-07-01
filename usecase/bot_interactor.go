package usecase

import "app/domain"

type BotInteractor struct {
	BotRepository BotRepository
}

func (interactor *BotInteractor) FetchAll() (bots domain.Bots, err error) {
	bots, err = interactor.BotRepository.FetchAll()
	return bots, err
}

func (interactor *BotInteractor) FetchOneById(identifier string) (bot domain.Bot, err error) {
	bot, err = interactor.BotRepository.FetchOneById(identifier)
	return bot, err
}
func (interactor *BotInteractor) FetchAllPublic() (bots domain.PublicBots, err error) {
	bots, err = interactor.BotRepository.FetchPublicAll()
	return bots, err
}

func (interactor *BotInteractor) FetchPublicOneById(identifier string) (bot domain.PublicBot, err error) {
	bot, err = interactor.BotRepository.FetchPublicOneById(identifier)
	return bot, err
}

func (interactor *BotInteractor) Create(b domain.Bot) (err error) {
	err = interactor.BotRepository.Create(b)
	return err
}
func (interactor *BotInteractor) Update(b *domain.Bot) (err error) {
	err = interactor.BotRepository.Update(b)
	return err
}
