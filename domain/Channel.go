package domain

import (
"time"

"gorm.io/gorm"
)

type Channel struct {
	ID             string    `gorm:"primaryKey"` // チャンネル ID
	Name           string    `gorm:"not null"`   // チャンネル名
	GuildID        string    `gorm:"not null;default:''"`
	BotID          string    `gorm:"not null;default:''"`
	Searchword     string    `gorm:"not null;default:''"`         // youtube で検索する際のハッシュタグ
	LastSearchedAt time.Time `gorm:"type:timestamp;default:null"` // 最後の検索 api を叩いた日時
	CreatedAt      time.Time `gorm:"not null;default:current_timestamp"`
	UpdatedAt      time.Time `gorm:"not null;default:current_timestamp;autoUpdateTime"`
	DeletedAt      gorm.DeletedAt `gorm:"type:timestamp;default:null"`
	Guild          Guild     `gorm:"foreignKey:GuildID"`
	Bot            Bot       `gorm:"foreignKey:BotID"`
}

type Channels []Channel
