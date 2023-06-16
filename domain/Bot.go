package domain

type Bot struct {
	ID          string `gorm:"primarykey"`
	Name        string `gorm:"not null"`
	Token       string `gorm:"unique;not null"`
	IsAvailable bool   `gorm:"default:false"`
}

type PublicBot struct {
	ID          string `gorm:"primarykey"`
	Name        string `gorm:"not null"`
	IsAvailable bool   `gorm:"default:false"`
}

type Bots []Bot
type PublicBots []PublicBot
