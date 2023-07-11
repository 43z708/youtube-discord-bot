package controllers

import (
	"app/domain"
	"app/interfaces/database"
	"app/usecase"
	"app/utilities"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"gorm.io/gorm"
)

type YoutubeController struct {
	BotInteractor       usecase.BotInteractor
	GuildInteractor     usecase.GuildInteractor
	ChannelInteractor   usecase.ChannelInteractor
	BlacklistInteractor usecase.BlacklistInteractor
}

type BotCron struct {
	CronJob *cron.Cron
	Running bool
}

var botCrons map[string]*BotCron = make(map[string]*BotCron)

func NewYoutubeController(sqlHandler *gorm.DB) *YoutubeController {
	return &YoutubeController{
		BlacklistInteractor: usecase.BlacklistInteractor{
			BlacklistRepository: &database.BlacklistRepository{
				SqlHandler: sqlHandler,
			},
		},
		ChannelInteractor: usecase.ChannelInteractor{
			ChannelRepository: &database.ChannelRepository{
				SqlHandler: sqlHandler,
			},
		},
		GuildInteractor: usecase.GuildInteractor{
			GuildRepository: &database.GuildRepository{
				SqlHandler: sqlHandler,
			},
		},
		BotInteractor: usecase.BotInteractor{
			BotRepository: &database.BotRepository{
				SqlHandler: sqlHandler,
			},
		},
	}
}

func (controller *YoutubeController) StartNotification(s *discordgo.Session, i *discordgo.InteractionCreate, newTimeInterval int64) {

	// DBにあるサーバー情報を取得
	guild, err := controller.GuildInteractor.FetchOneById(i.GuildID)

	if err != nil {
		utilities.InteractionReply(s, i, "データの取得に失敗しました。:"+err.Error())
		return
	}
	if guild.AdminChannelID != i.ChannelID {
		utilities.InteractionReply(s, i, fmt.Sprintf("このコマンドは <#%s> でのみ許可されたコマンドです。", guild.AdminChannelID))
		return
	}

	// bot稼働時にDBから取得した情報から更新されている可能性があるためDB再取得
	bot, err := controller.BotInteractor.FetchOneById(s.State.User.ID)

	// チャンネルの設定があるかどうかをチェック
	channels, err := controller.ChannelInteractor.FetchAllByGuildID(i.GuildID)
	if err != nil && i != nil {
		utilities.InteractionReply(s, i, "チャンネルを設定してからでないと通知はできません。")
		return
	}

	if newTimeInterval > 0 {
		// 新しい時間間隔を指定された場合
		if botCrons[s.State.User.ID] != nil && botCrons[s.State.User.ID].Running && bot.TimeInterval == newTimeInterval {
			// 新しい時間間隔が既存の時間間隔と一致している場合、エラーを返す
			utilities.InteractionNonEphemeralReply(s, i, "すでに"+strconv.FormatInt(bot.TimeInterval, 10)+"時間おきにyoutube通知設定されています。")
			return
		} else if botCrons[s.State.User.ID] != nil && botCrons[s.State.User.ID].Running && bot.TimeInterval != newTimeInterval {
			// 新しい時間間隔が既存と異なる場合、いったんcronをとめてDBを更新
			botCrons[s.State.User.ID].CronJob.Stop()

			if err != nil {
				utilities.InteractionReply(s, i, "データの取得に失敗しました。:"+err.Error())
				return
			}
			bot.TimeInterval = newTimeInterval

			err = controller.BotInteractor.Update(&bot)
			if err != nil {
				utilities.InteractionReply(s, i, "データの更新に失敗しました。:"+err.Error())
				return
			}
		}
	} else if botCrons[s.State.User.ID] != nil && botCrons[s.State.User.ID].Running {
		// 時間指定はしなかったけど、すでに稼働している場合
		utilities.InteractionNonEphemeralReply(s, i, "すでに"+strconv.FormatInt(bot.TimeInterval, 10)+"時間おきにyoutube通知設定されています。")
		return
	}

	// YouTubeクライアントの作成
	youtubeAPIKey := os.Getenv("YOUTUBE_API_KEY")

	ctx := context.Background()
	ytSvc, err := youtube.NewService(ctx, option.WithAPIKey(youtubeAPIKey))
	if err != nil {
		log.Fatalf("Error creating YouTube client: %s", err.Error())
	}

	c := cron.New()
	_, err = c.AddFunc("0 */"+strconv.FormatInt(bot.TimeInterval, 10)+" * * *", func() {
		controller.PostLatestYoutubeVideo(ytSvc, s, guild)
	})
	if err != nil {
		utilities.InteractionReply(s, i, "通知設定に失敗しました。:"+err.Error())
		return
	}
	c.Start()
	botCrons[s.State.User.ID] = &BotCron{CronJob: c, Running: true}

	// 通知開始
	var channelIds []string
	for _, channel := range channels {
		channelIds = append(channelIds, "キーワード : "+channel.Searchword+" => <#"+channel.ID+">")
	}
	channelIdLinks := strings.Join(channelIds, "\n")
	utilities.InteractionNonEphemeralReply(s, i, "youtube動画の通知設定を"+strconv.FormatInt(bot.TimeInterval, 10)+"時間おきにして稼働開始しました。\n配信対象チャンネルは以下です。\n"+channelIdLinks)
}

func (controller *YoutubeController) PostLatestYoutubeVideo(ytSvc *youtube.Service, s *discordgo.Session, guild domain.Guild) {

	// blacklistを取得
	blacklists, err := controller.BlacklistInteractor.FetchAllByGuildID(guild.ID)
	if err != nil {
		s.ChannelMessageSend(guild.AdminChannelID, "Blacklistの取得に失敗しました。:"+err.Error())
		return
	}

	channels, err := controller.ChannelInteractor.FetchAllByGuildID(guild.ID)

	if err != nil {
		s.ChannelMessageSend(guild.AdminChannelID, "チャンネル情報の取得に失敗しました。:"+err.Error())
		return
	}

	// チャンネルごとにyoutubeのAPIを叩く
	for _, channel := range channels {
		if channel.Searchword == "" {
			s.ChannelMessageSend(guild.AdminChannelID, fmt.Sprintf("<#%s> に検索キーワードの設定がされていないため取得に失敗しました。", channel.ID))
			return
		}

		var nextPageToken string

		// youtubeのAPIは一度につき50件までなのでpagenationで対応
		for {
			call := ytSvc.Search.List([]string{"id", "Snippet"})
			call.Q(channel.Searchword).MaxResults(50).Type("video").PublishedAfter(channel.LastSearchedAt.Format(time.RFC3339))
			if nextPageToken != "" {
				call.PageToken(nextPageToken)
			}
			response, err := call.Do()
			if err != nil {
				s.ChannelMessageSend(guild.AdminChannelID, "Youtube APIでのデータ取得に失敗しました。:"+err.Error())
				return
			}
			log.Println(&response.Items)
			for _, video := range response.Items {
				channelId := video.Snippet.ChannelId
				log.Println(channelId)
				var shouldBreak = false
				for _, blacklist := range blacklists {
					if channelId == blacklist.Distributor {
						shouldBreak = true
						videoURL := fmt.Sprintf("https://www.youtube.com/watch?v=%s", video.Id.VideoId)
						message := fmt.Sprintf("この動画はブラックリストに含まれているチャンネルのため通知を拒否しました: %s", videoURL)
						_, err := s.ChannelMessageSend(guild.AdminChannelID, message)
						if err != nil {
							s.ChannelMessageSend(guild.AdminChannelID, "動画URLの送信に失敗しました。:"+err.Error())
						}
						break
					}
				}
				if shouldBreak {
					// ブラックリストに含まれたチャンネルの動画の場合は投稿しない
					break
				}
				videoURL := fmt.Sprintf("https://www.youtube.com/watch?v=%s", video.Id.VideoId)
				message := fmt.Sprintf("New video: %s", videoURL)
				_, err := s.ChannelMessageSend(channel.ID, message)
				if err != nil {
					s.ChannelMessageSend(guild.AdminChannelID, "動画URLの送信に失敗しました。:"+err.Error())
				}
			}

			nextPageToken = response.NextPageToken
			if nextPageToken == "" {
				break
			}
			channel.LastSearchedAt = time.Now()
			err = controller.ChannelInteractor.Update(&channel)
			if err != nil {
				s.ChannelMessageSend(guild.AdminChannelID, "DBの更新に失敗しました。:"+err.Error())
			}
		}

	}

	var channelIds []string
	for _, channel := range channels {
		channelIds = append(channelIds, "キーワード : "+channel.Searchword+" => <#"+channel.ID+">")
	}
	channelIdLinks := strings.Join(channelIds, "\n")

	s.ChannelMessageSend(guild.AdminChannelID, "以下のチャンネルの動画投稿通知が完了しました。\n"+channelIdLinks)
}

func (controller *YoutubeController) StopNotification(s *discordgo.Session, i *discordgo.InteractionCreate) {

	// DBにあるサーバー情報を取得
	guild, err := controller.GuildInteractor.FetchOneById(i.GuildID)

	if err != nil {
		utilities.InteractionReply(s, i, "データの取得に失敗しました。:"+err.Error())
		return
	}
	if guild.AdminChannelID != i.ChannelID {
		utilities.InteractionReply(s, i, fmt.Sprintf("このコマンドは <#%s> でのみ許可されたコマンドです。", guild.AdminChannelID))
		return
	}
	if botCrons[s.State.User.ID] == nil || !botCrons[s.State.User.ID].Running {
		utilities.InteractionNonEphemeralReply(s, i, "すでに稼働しているyoutube動画通知はありません。")
		return
	}

	botCrons[s.State.User.ID].CronJob.Stop()
	botCrons[s.State.User.ID].Running = false
	utilities.InteractionNonEphemeralReply(s, i, "youtube動画通知の稼働を停止しました。")
}
