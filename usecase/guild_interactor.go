package usecase

import "app/domain"

type GuildInteractor struct {
	GuildRepository GuildRepository
}

func (interactor *GuildInteractor) FetchAll() (guilds domain.Guilds, err error) {
	guilds, err = interactor.GuildRepository.FetchAll()
	return guilds, err
}

func (interactor *GuildInteractor) FetchOneById(identifier string) (guild domain.Guild, err error) {
	guild, err = interactor.GuildRepository.FetchOneById(identifier)
	return guild, err
}
func (interactor *GuildInteractor) FetchPublicAllByBotID(identifier string) (guilds domain.PublicGuilds, err error) {
	guilds, err = interactor.GuildRepository.FetchPublicAllByBotID(identifier)
	return guilds, err
}

func (interactor *GuildInteractor) FetchPublicOneById(identifier string) (guild domain.PublicGuild, err error) {
	guild, err = interactor.GuildRepository.FetchPublicOneById(identifier)
	return guild, err
}

func (interactor *GuildInteractor) Create(b domain.Guild) (err error) {
	_, err = interactor.GuildRepository.Create(b)
	return err
}

func (interactor *GuildInteractor) Update(b domain.Guild) (err error) {
	_, err = interactor.GuildRepository.Update(b)
	return err
}
func (interactor *GuildInteractor) Delete(identifier string) (err error) {
	_, err = interactor.GuildRepository.Delete(identifier)
	return err
}
