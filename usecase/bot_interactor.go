package usecase

import "app/domain"

type BotInteractor struct {
	BotRepository BotRepository
}

func (interactor *BotInteractor) Add(b domain.Bot) (err error) {
	_, err = interactor.BotRepository.Store(b)
	return err
}

func (interactor *BotInteractor) Bots() (bots domain.Bots, err error) {
	bots, err = interactor.BotRepository.FindAll()
	return bots, err
}

func (interactor *BotInteractor) BotById(identifier int) (bot domain.Bot, err error) {
	bot, err = interactor.BotRepository.FindById(identifier)
	return bot, err
}
func (interactor *BotInteractor) PublicBots() (bots domain.PublicBots, err error) {
	bots, err = interactor.BotRepository.FindPublicAll()
	return bots, err
}

func (interactor *BotInteractor) PublicBotById(identifier int) (bot domain.PublicBot, err error) {
	bot, err = interactor.BotRepository.FindPublicById(identifier)
	return bot, err
}
