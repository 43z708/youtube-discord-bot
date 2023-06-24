package domain

import "time"

type Channel struct {
	ID             string    `gorm:"primaryKey"`
	Name           string    `gorm:"not null"`
	GuildID        string    `gorm:"not null;default:''"`
	BotID          string    `gorm:"not null;default:''"`
	Searchword     string    `gorm:"not null;default:''"`
	LastSearchedAt time.Time `gorm:"type:timestamp;default:null"`
	CreatedAt      time.Time `gorm:"not null;default:current_timestamp"`
	UpdatedAt      time.Time `gorm:"not null;default:current_timestamp;autoUpdateTime"`
	DeletedAt      *time.Time
	Guild          Guild `gorm:"foreignKey:GuildID"`
	Bot            Bot   `gorm:"foreignKey:BotID"`
}

type Channels []Channel
