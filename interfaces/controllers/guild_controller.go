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
			s.ChannelMessageSend(guild.AdminChannelID, "まずは/register-apikeyコマンドでyoutubeのapikeyを登録してください。")
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
			utilities.InteractionReply(s, i, "Botの準備が完了しました。\"/create-channel\"コマンドでbotが投稿するチャンネルと検索ワードを設定・作成すれば設定完了です。")
			return
		}

	}
}

func (controller *GuildController) Delete(s *discordgo.Session, c *discordgo.ChannelDelete) {
}
