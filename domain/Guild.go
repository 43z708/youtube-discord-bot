package domain

type Guild struct {
	ID            string `gorm:"primaryKey"`
	Name          string `gorm:"not null"`
	YoutubeApiKey string `gorm:"unique;default:null"`
	BotID         string `gorm:"default:null"`
	Bot           Bot    `gorm:"foreignKey:BotID"`
}
type PublicGuild struct {
	ID    string `gorm:"primaryKey"`
	Name  string `gorm:"not null"`
	BotID string `gorm:"default:null"`
	Bot   Bot    `gorm:"foreignKey:BotID"`
}

type Guilds []Guild
type PublicGuilds []PublicGuild
