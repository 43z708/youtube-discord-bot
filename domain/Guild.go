package domain

type Guild struct {
	ID    string `gorm:"primaryKey"`
	Name  string `gorm:"not null"`
	BotID string `gorm:"default:null"`
	Bot   Bot    `gorm:"foreignKey:BotID"`
}
