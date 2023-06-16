package infrastructure

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func Youtube() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading dotenv: %s", err.Error())
	}

	var (
		ctx   context.Context
		ytSvc *youtube.Service
	)

	// YouTubeクライアントの作成
	youtubeAPIKey := os.Getenv("YOUTUBE_API_KEY")

	ctx = context.Background()
	ytSvc, err = youtube.NewService(ctx, option.WithAPIKey(youtubeAPIKey))
	if err != nil {
		log.Fatalf("Error creating YouTube client: %s", err.Error())
	}

	// cronジョブの作成
	c := cron.New()

	// 1分(動作確認のため仮)ごとに動画を投稿するジョブを追加
	// 3時間おきの場合、 0 */3 * * *
	_, err = c.AddFunc("*/1 * * * *", func() {
		postLatestYouTubeVideo(ytSvc, "1110531236151173150", "#test")
	})
	if err != nil {
		log.Fatalf("Error adding cron job: %s", err.Error())
	}

	// cronジョブの実行
	c.Start()

	// シグナルを受け取るためのチャネルを作成
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop // シグナルを受け取るまでブロックされる

	// cronジョブの停止
	c.Stop()
}

func postLatestYouTubeVideo(ytSvc *youtube.Service, channelID string, searchQuery string) {

	call := ytSvc.Search.List([]string{"id"})
	call.Q(searchQuery).MaxResults(1).Type("video")
	response, err := call.Do()
	if err != nil {
		log.Printf("Error searching YouTube videos: %s", err.Error())
		return
	}

	if len(response.Items) > 0 {
		video := response.Items[0]
		videoURL := fmt.Sprintf("https://www.youtube.com/watch?v=%s", video.Id.VideoId)
		message := fmt.Sprintf("New video: %s", videoURL)
		_, err := dg.ChannelMessageSend(channelID, message)
		if err != nil {
			log.Printf("Error sending message to Discord channel: %s", err.Error())
		}
	}
}
