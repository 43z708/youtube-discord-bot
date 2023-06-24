package domain

type Guild struct {
	ID            string `gorm:"primaryKey"`          //サーバーの id
	Name          string `gorm:"not null"`            //サーバー名
	YoutubeApiKey string `gorm:"unique;default:null"` //  youtubeのapi key
	CategoryName  string `gorm:"default:null"`        // 投稿チャンネルのカテゴリ名
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
