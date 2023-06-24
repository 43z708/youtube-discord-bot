package domain

type Bot struct {
	ID          string `gorm:"primarykey"`      // bot の id (application id)
	Name        string `gorm:"not null"`        // bot の名前
	Token       string `gorm:"unique;not null"` // bot の api トークン
	IsAvailable bool   `gorm:"default:false"`   // 有効な bot かどうか
}

type PublicBot struct {
	ID          string `gorm:"primarykey"`
	Name        string `gorm:"not null"`
	IsAvailable bool   `gorm:"default:false"`
}

type Bots []Bot
type PublicBots []PublicBot
