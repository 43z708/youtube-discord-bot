package main

import (
	"app/infrastructure"
	"context"

	"google.golang.org/api/youtube/v3"
)

var (
	ctx   context.Context
	ytSvc *youtube.Service
)

func main() {

	Init := infrastructure.Init()
	infrastructure.Router(Init)
	infrastructure.Discord(Init)

	infrastructure.Youtube()

	// // Discordセッションの作成
	// controller.CreateSession()
	// // Discordセッションを開始
	// controller.StartSession()
	// // YouTubeクライアントの作成
	// controller.CreateClient()
	// // cronの実行
	// controller.StartCron()

	// log.Println("Bot stopped.")
}
