package domain

import "time"

type Channel struct {
	ID         string    `gorm:"primaryKey"`
	Name       string    `gorm:"not null"`
	GuildID    string    `gorm:"not null;default:''"`
	BotID      string    `gorm:"not null;default:''"`
	SearchWord string    `gorm:"not null;default:''"`
	SearchedAt time.Time `gorm:"type:timestamp;not null;default:current_timestamp"`
	CreatedAt  time.Time `gorm:"not null;default:current_timestamp"`
	UpdatedAt  time.Time `gorm:"not null;default:current_timestamp;autoUpdateTime"`
	DeletedAt  *time.Time
	Guild      Guild `gorm:"foreignKey:GuildID"`
	Bot        Bot   `gorm:"foreignKey:BotID"`
}
