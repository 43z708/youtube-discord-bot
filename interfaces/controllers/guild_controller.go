package controllers

import (
	"app/domain"
	"app/interfaces/database"
	"app/usecase"
	"app/utilities"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type GuildController struct {
	GuildInteractor usecase.GuildInteractor
}

func NewGuildController(sqlHandler *gorm.DB) *GuildController {
	return &GuildController{
		GuildInteractor: usecase.GuildInteractor{
			GuildRepository: &database.GuildRepository{
				SqlHandler: sqlHandler,
			},
		},
	}
}

func (controller *GuildController) FetchPublicOneById(id string) domain.PublicGuild {
	guild, err := controller.GuildInteractor.FetchPublicOneById(id)

	if err != nil {
		log.Fatalf("Error getting guilds data: %s", err.Error())
	}
	return guild
}

func (controller *GuildController) FetchPublicAllByBotID(id string) domain.PublicGuilds {
	guilds, err := controller.GuildInteractor.FetchPublicAllByBotID(id)
	if err != nil {
		log.Fatalf("Error getting channels data: %s", err.Error())
	}

	return guilds
}

func (controller *GuildController) FetchOneById(id string) domain.Guild {
	guild, err := controller.GuildInteractor.FetchOneById(id)

	if err != nil {
		log.Fatalf("Error getting guilds data: %s", err.Error())
	}
	return guild
}

func (controller *GuildController) FetchAll() domain.Guilds {
	guilds, err := controller.GuildInteractor.FetchAll()
	if err != nil {
		log.Fatalf("Error getting guilds data: %s", err.Error())
	}

	return guilds
}

func (controller *GuildController) Create(s *discordgo.Session, e *discordgo.GuildCreate) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading dotenv: %s", err.Error())
	}

	// コマンドの登録
	utilities.RegisterCommand(s, e)

	// DBにあるサーバー情報を取得
	guild, err := controller.GuildInteractor.FetchOneById(e.Guild.ID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// DBにサーバー情報がない場合、チャンネル等作成して保存

			category, err := s.GuildChannelCreateComplex(e.Guild.ID, discordgo.GuildChannelCreateData{
				Name: os.Getenv("CATEGORY_NAME"),
				Type: discordgo.ChannelTypeGuildCategory,
			})
			if err != nil {
				log.Println("Error creating category:", err)
				return
			}

			channel, err := s.GuildChannelCreateComplex(e.Guild.ID, discordgo.GuildChannelCreateData{
				Name:     "admin-channel",
				Type:     discordgo.ChannelTypeGuildText,
				ParentID: category.ID,
				PermissionOverwrites: []*discordgo.PermissionOverwrite{
					{
						ID:   e.Guild.ID, // This is the ID of the @everyone role
						Type: discordgo.PermissionOverwriteTypeRole,
						Deny: discordgo.PermissionViewChannel, // Deny view permission for everyone
					},
					{
						ID:   e.Guild.OwnerID, // This is the ID of the admin
						Type: discordgo.PermissionOverwriteTypeMember,
						Allow: discordgo.PermissionAdministrator |
							discordgo.PermissionViewChannel |
							discordgo.PermissionSendMessages |
							discordgo.PermissionReadMessageHistory,
					},
					{
						ID:   s.State.User.ID, // This is the ID of the bot
						Type: discordgo.PermissionOverwriteTypeMember,
						Allow: discordgo.PermissionViewChannel |
							discordgo.PermissionSendMessages |
							discordgo.PermissionReadMessageHistory,
					},
				},
			})
			if err != nil {
				s.ChannelMessageSend(channel.ID, "チャンネルの作成に失敗しました。:"+err.Error())
				return
			}
			guild := domain.Guild{
				ID:             e.Guild.ID,
				Name:           e.Guild.Name,
				CategoryID:     category.ID,
				BotID:          s.State.User.ID,
				AdminChannelID: channel.ID,
			}

			controller.GuildInteractor.Create(guild)
			if err != nil {
				s.ChannelMessageSend(channel.ID, "DBの保存に失敗しました。:"+err.Error())
				return
			}
			s.ChannelMessageSend(guild.AdminChannelID, "まずは`/register-apikey`コマンドでyoutubeのapikeyを登録してください。`/help`でbotの使い方を見ることができます。")
			return
		}
	} else {
		// guildデータがDBに存在した場合（botがキックされたケースを想定）
		categoryID := ""
		if guild.CategoryID == "" {
			category, err := s.GuildChannelCreateComplex(e.Guild.ID, discordgo.GuildChannelCreateData{
				Name: os.Getenv("CATEGORY_NAME"),
				Type: discordgo.ChannelTypeGuildCategory,
			})
			if err != nil {
				log.Println("Error creating category:", err)
				return
			}
			categoryID = category.ID
			guild.CategoryID = categoryID
		} else {
			categoryID = guild.CategoryID
		}

		if guild.AdminChannelID == "" {
			channel, err := s.GuildChannelCreateComplex(e.Guild.ID, discordgo.GuildChannelCreateData{
				Name:     "admin-channel",
				Type:     discordgo.ChannelTypeGuildText,
				ParentID: categoryID,
				PermissionOverwrites: []*discordgo.PermissionOverwrite{
					{
						ID:   e.Guild.ID, // This is the ID of the @everyone role
						Type: discordgo.PermissionOverwriteTypeRole,
						Deny: discordgo.PermissionViewChannel, // Deny view permission for everyone
					},
					{
						ID:   e.Guild.OwnerID, // This is the ID of the admin
						Type: discordgo.PermissionOverwriteTypeMember,
						Allow: discordgo.PermissionAdministrator |
							discordgo.PermissionViewChannel |
							discordgo.PermissionSendMessages |
							discordgo.PermissionReadMessageHistory,
					},
					{
						ID:   s.State.User.ID, // This is the ID of the bot
						Type: discordgo.PermissionOverwriteTypeMember,
						Allow: discordgo.PermissionViewChannel |
							discordgo.PermissionSendMessages |
							discordgo.PermissionReadMessageHistory,
					},
				},
			})
			if err != nil {
				s.ChannelMessageSend(channel.ID, "チャンネルの作成に失敗しました。:"+err.Error())
				return
			}
			guild.AdminChannelID = channel.ID
		}

		controller.GuildInteractor.Update(&guild)

	}
}

func (controller *GuildController) Update(s *discordgo.Session, i *discordgo.InteractionCreate) {

	// コマンドの場合（最初に発火）
	if i.Type == discordgo.InteractionApplicationCommand {

		// スラッシュコマンドのデータを取得する
		command := i.ApplicationCommandData()

		// /register-apikey コマンド以外は無視する
		if command.Name == "register-apikey" {

			// DBにあるサーバー情報を取得
			guild, err := controller.GuildInteractor.FetchOneById(i.GuildID)

			if guild.AdminChannelID != i.ChannelID {
				utilities.InteractionReply(s, i, fmt.Sprintf("このコマンドは <#%s> でのみ許可されたコマンドです。", guild.AdminChannelID))
				return
			}

			// youtube api keyを取得
			apiKey := command.Options[0].StringValue()

			if err != nil {
				utilities.InteractionReply(s, i, "データの取得に失敗しました。:"+err.Error())
				return
			}
			guild.YoutubeApiKey = apiKey

			err = controller.GuildInteractor.Update(&guild)

			if err != nil {
				utilities.InteractionReply(s, i, "データの更新に失敗しました。:"+err.Error())
				return
			}

			// カテゴリ名がすでに登録されている場合は終了
			utilities.InteractionReply(s, i, "Botの準備が完了しました。`/create-channel`コマンドでbotが投稿するチャンネルと検索ワードを必要な分だけ設定し、`/start-notification`コマンドでyoutube動画の通知が稼働します。\nコマンド一覧は`/help`で出力されます。")
			return
		}

	}
}

func (controller *GuildController) Delete(s *discordgo.Session, c *discordgo.ChannelDelete) {
}

func (controller *GuildController) Help(s *discordgo.Session, i *discordgo.InteractionCreate) {

	// コマンドの場合（最初に発火）
	if i.Type == discordgo.InteractionApplicationCommand {

		// スラッシュコマンドのデータを取得する
		command := i.ApplicationCommandData()

		// /register-apikey コマンド以外は無視する
		if command.Name == "help" {

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

			// コマンド一覧
			embed := &discordgo.MessageEmbed{
				Author: &discordgo.MessageEmbedAuthor{},
				Color:  0x00ff00, // Green
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "`/register-apikey`コマンド",
						Value:  "Youtube動画投稿のためのapi keyの登録コマンド。",
						Inline: false,
					},
					{
						Name:   "`/create-channel`コマンド",
						Value:  "指定した検索キーワードのyoutube動画を投稿するチャンネルを作るコマンド。\nchannel_nameにチャンネル名、search_wordsに検索キーワードを入力。\n検索キーワードは、\"Edit Channel\"→\"Overview\"→\"CHANNEL TOPIC\"にて変更可能。",
						Inline: false,
					},
					{
						Name:   "`/start-notification`コマンド",
						Value:  "time_intervalに指定した時間間隔（単位：時間）でyoutubeの検索結果の通知を開始できるコマンド。\ntime_intervalに何も指定しない場合、defaultで3時間の設定となる。\n入力できる数値は1から23まで。",
						Inline: false,
					},
					{
						Name:   "`/stop-notification`コマンド",
						Value:  "youtubeの検索結果の通知を停止できるコマンド。",
						Inline: false,
					},
					{
						Name:   "`/add-blacklist`コマンド",
						Value:  "ブラックリストに含めたいチャンネルを設定できる。チャンネルIDを入力することで設定可能。" + utilities.ExplainGetYoutubeChannelID,
						Inline: false,
					},
					{
						Name:   "`/remove-blacklist`コマンド",
						Value:  "ブラックリストに指定したチャンネルを削除できる。チャンネルIDを入力することで設定可能。",
						Inline: false,
					},
					{
						Name:   "`/get-blacklist`コマンド",
						Value:  "ブラックリストに指定したチャンネル一覧を出力。",
						Inline: false,
					},
				},
				Title: "コマンド一覧",
			}

			data := &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{embed},
				Flags:  discordgo.MessageFlagsEphemeral,
			}

			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: data,
			})

			return
		}

	}
}
