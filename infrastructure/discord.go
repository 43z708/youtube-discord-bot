package infrastructure

import (
	"app/domain"
	"app/interfaces/controllers"
	"log"

	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

var (
	dg  *discordgo.Session
	err error
)

func Discord(Init *gorm.DB) {

	BotController := controllers.NewBotController(Init)

	bots := BotController.FetchAll()

	CreateSession(bots, Init)

	StartSession()

}

// Discordセッションの作成
func CreateSession(bots domain.Bots, Init *gorm.DB) {
	ChannelController := controllers.NewChannelController(Init)
	for _, bot := range bots {
		discordToken := bot.Token
		dg, err = discordgo.New("Bot " + discordToken)
		if err != nil {
			log.Fatalf("Error creating Discord session: %s", err.Error())
		}

		dg.AddHandler(ChannelController.Create)
		dg.AddHandler(ChannelController.Update)
		dg.AddHandler(ChannelController.Delete)
	}
}

// Discordセッションを開始
func StartSession() {
	err = dg.Open()
	if err != nil {
		log.Fatalf("Error opening Discord session: %s", err.Error())
	}
	defer dg.Close()

	log.Println("Discord bot is running!")
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// メッセージが受信されたときに呼ばれるハンドラ
	// ここで必要な処理を追加してください
}
