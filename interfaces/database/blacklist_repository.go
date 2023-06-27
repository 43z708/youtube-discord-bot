package database

import (
	"app/domain"
	"strconv"

	"gorm.io/gorm"
)

type BlacklistRepository struct {
	SqlHandler *gorm.DB
}

func (repo *BlacklistRepository) FetchOneById(i string) (domain.Blacklist, error) {
	var blacklist domain.Blacklist
	id, _ := strconv.Atoi(i)
	result := repo.SqlHandler.First(&blacklist, id)
	return blacklist, result.Error
}

func (repo *BlacklistRepository) FetchOneByDistributor(Distributor string) (domain.Blacklist, error) {
	var blacklist domain.Blacklist
	result := repo.SqlHandler.Where("distributor = ?", Distributor).First(&blacklist)
	return blacklist, result.Error
}

func (repo *BlacklistRepository) FetchAllByGuildID(guildID string) (domain.Blacklists, error) {
	blacklists := make([]domain.Blacklist, 0)
	result := repo.SqlHandler.Where("guild_id = ?", guildID).Find(&blacklists)
	return blacklists, result.Error
}

func (repo *BlacklistRepository) Create(b domain.Blacklist) error {

	result := repo.SqlHandler.Create(&b)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *BlacklistRepository) Update(blacklist *domain.Blacklist) error {
	result := repo.SqlHandler.Save(&blacklist)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *BlacklistRepository) Delete(id int) error {
	var blacklist domain.Blacklist
	result := repo.SqlHandler.First(&blacklist, id)
	result = repo.SqlHandler.Delete(&blacklist)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
