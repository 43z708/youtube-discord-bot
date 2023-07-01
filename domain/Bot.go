package domain

type Bot struct {
	ID           string `gorm:"primarykey"`      // bot の id (application id)
	Name         string `gorm:"not null"`        // bot の名前
	TimeInterval int64  `gorm:"default:3"`       // botの投稿頻度(デフォルト3時間おき)
	Token        string `gorm:"unique;not null"` // bot の api トークン
	IsAvailable  bool   `gorm:"default:false"`   // 有効な bot かどうか
}

type PublicBot struct {
	ID           string `gorm:"primarykey"`
	Name         string `gorm:"not null"`
	TimeInterval int64  `gorm:"default:3"`
	IsAvailable  bool   `gorm:"default:false"`
}

type Bots []Bot
type PublicBots []PublicBot
