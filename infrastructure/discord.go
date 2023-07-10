package infrastructure

import (
	"app/domain"
	"app/interfaces/controllers"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
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

}

// Discordセッションの作成
func CreateSession(bots domain.Bots, Init *gorm.DB) {
	var wg sync.WaitGroup

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading dotenv: %s", err.Error())
	}

	appEnv := os.Getenv("APP_ENV")

	ChannelController := controllers.NewChannelController(Init)
	GuildController := controllers.NewGuildController(Init)
	BlacklistController := controllers.NewBlacklistController(Init)
	YoutubeController := controllers.NewYoutubeController(Init)

	for _, bot := range bots {
		if !bot.IsAvailable {
			break
		}
		wg.Add(1) // WaitGroupのカウンタを増やす

		log.Println(bot.Name + " init")
		discordToken := bot.Token
		dg, err = discordgo.New("Bot " + discordToken)
		if err != nil {
			log.Fatalf("Error creating Discord session: %s", err.Error())
		}
		go func(bot domain.Bot) {

			// create-channelコマンドの処理
			dg.AddHandler(ChannelController.Create)
			// DBに保存しているチャンネル情報が変更されたときの処理
			dg.AddHandler(ChannelController.Update)
			// DBに保存されているチャンネルがdiscord側で削除されたときにDBも削除する処理
			dg.AddHandler(ChannelController.Delete)
			// bot招待時にカテゴリとadmin-channelを作成しDBに保存する処理
			dg.AddHandler(GuildController.Create)
			// register-apikeyコマンドの処理
			dg.AddHandler(GuildController.Update)
			// get-blacklistコマンド
			dg.AddHandler(BlacklistController.FetchBlacklist)
			// add-blacklistコマンド
			dg.AddHandler(BlacklistController.Create)
			// remove-blacklistコマンド
			dg.AddHandler(BlacklistController.Delete)
			// helpコマンド
			dg.AddHandler(GuildController.Help)
			// /start-notificationと/stop-notificationコマンド
			dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {

				// コマンドの場合（最初に発火）
				if i.Type == discordgo.InteractionApplicationCommand {

					// スラッシュコマンドのデータを取得する
					command := i.ApplicationCommandData()

					// /start-notification コマンド以外は無視する
					if command.Name == "start-notification" {
						var time_interval int64 = 0
						if len(command.Options) > 0 {
							time_interval = command.Options[0].IntValue()
						}
						YoutubeController.StartNotification(s, i, time_interval)
					} else if command.Name == "stop-notification" {
						YoutubeController.StopNotification(s, i)
					}
				}

			})

			if appEnv != "local" {
				guilds := GuildController.FetchPublicAllByBotID(dg.State.User.ID)
				for _, guild := range guilds {
					if guild.AdminChannelID != "" {
						dg.ChannelMessageSend(guild.AdminChannelID, "botが稼働しました。現在通知が止まっています。再稼働させるには`/start-notification`コマンドを入力してください。")
					}
				}
			}

			wg.Done() // 処理が終了したらWaitGroupのカウンタを減らす
		}(bot)

		err = dg.Open()
		if err != nil {
			log.Fatalf("Error opening Discord session: %s", err.Error())
		}
		defer dg.Close()

		log.Println("Discord bot is running!")
	}

	// Wait here until CTRL-C or other term signal is received.
	log.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	wg.Wait() // WaitGroupのカウンタが0になる（すべてのgoroutineが終了する）まで待つ

}
