package domain

type Blacklist struct {
	ID          int `gorm:"primaryKey;autoIncrement"`
	Distributor string
	BotID       string
	GuildID     string
	Guild       Guild `gorm:"foreignKey:GuildID"`
	Bot         Bot   `gorm:"foreignKey:BotID"`
}

type Blacklists = []Blacklist
