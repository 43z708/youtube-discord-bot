package usecase

import "app/domain"

type ChannelInteractor struct {
	ChannelRepository ChannelRepository
}

func (interactor *ChannelInteractor) FetchOneById(identifier string) (channel domain.Channel, err error) {
	channel, err = interactor.ChannelRepository.FetchOneById(identifier)
	return channel, err
}

func (interactor *ChannelInteractor) FetchAllByBotID(identifier string) (channels domain.Channels, err error) {
	channels, err = interactor.ChannelRepository.FetchAllByBotID(identifier)
	return channels, err
}

func (interactor *ChannelInteractor) FetchAllByGuildID(identifier string) (channels domain.Channels, err error) {
	channels, err = interactor.ChannelRepository.FetchAllByGuildID(identifier)
	return channels, err
}

func (interactor *ChannelInteractor) Create(b domain.Channel) (err error) {
	err = interactor.ChannelRepository.Create(b)
	return err
}

func (interactor *ChannelInteractor) Update(b *domain.Channel) (err error) {
	err = interactor.ChannelRepository.Update(b)
	return err
}
func (interactor *ChannelInteractor) Delete(identifier string) (err error) {
	err = interactor.ChannelRepository.Delete(identifier)
	return err
}
