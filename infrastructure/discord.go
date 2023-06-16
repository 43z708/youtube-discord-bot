package infrastructure

import (
	"app/domain"
	"app/interfaces/controllers"
	"log"

	"github.com/bwmarrin/discordgo"
)

var (
	dg  *discordgo.Session
	err error
)

func Discord() {

	BotController := controllers.NewBotController(Init())

	bots := BotController.Index()

	CreateSession(bots)

	StartSession()

}

// Discordセッションの作成
func CreateSession(bots domain.Bots) {
	for _, bot := range bots {
		discordToken := bot.Token
		dg, err = discordgo.New("Bot " + discordToken)
		if err != nil {
			log.Fatalf("Error creating Discord session: %s", err.Error())
		}

		dg.AddHandler(messageCreate)
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
