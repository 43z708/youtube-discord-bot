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
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type ChannelController struct {
	ChannelInteractor usecase.ChannelInteractor
	GuildInteractor   usecase.GuildInteractor
}

func NewChannelController(sqlHandler *gorm.DB) *ChannelController {
	return &ChannelController{
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
	}
}

func (controller *ChannelController) FetchOneById(id string) domain.Channel {
	channel, err := controller.ChannelInteractor.FetchOneById(id)

	if err != nil {
		log.Fatalf("Error getting channels data: %s", err.Error())
	}
	return channel
}

func (controller *ChannelController) FetchAllByBotID(id string) domain.Channels {
	channels, err := controller.ChannelInteractor.FetchAllByBotID(id)
	if err != nil {
		log.Fatalf("Error getting channels data: %s", err.Error())
	}

	return channels
}

func (controller *ChannelController) FetchAllByGuildID(id string) domain.Channels {
	channels, err := controller.ChannelInteractor.FetchAllByGuildID(id)
	if err != nil {
		log.Fatalf("Error getting channels data: %s", err.Error())
	}

	return channels
}

func (controller *ChannelController) Create(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionApplicationCommand {
		// スラッシュコマンドのデータを取得する
		command := i.ApplicationCommandData()

		// /create-channel コマンド以外は無視する
		if command.Name == "create-channel" {

			// DBにあるサーバー情報を取得
			guild, err := controller.GuildInteractor.FetchOneById(i.GuildID)

			if err != nil {
				utilities.InteractionReply(s, i, "DBからのサーバー情報の取得に失敗しました。:"+err.Error())
				return
			}

			if guild.AdminChannelID != i.ChannelID {
				utilities.InteractionReply(s, i, fmt.Sprintf("このコマンドは <#%s> でのみ許可されたコマンドです。", guild.AdminChannelID))
				return
			}

			// チャンネル名とタグのパラメータを取得する
			channelName := command.Options[0].StringValue()
			tag := command.Options[1].StringValue()

			// youtubeリンク投稿カテゴリの存在チェック
			err = godotenv.Load()
			if err != nil {
				log.Fatalf("Error loading dotenv: %s", err.Error())
			}
			// 存在しない場合カテゴリを作成してdbの更新
			if guild.CategoryID == "" {
				newCategory, err := s.GuildChannelCreateComplex(i.GuildID, discordgo.GuildChannelCreateData{
					Name:     os.Getenv("CATEGORY_NAME"),
					Type:     discordgo.ChannelTypeGuildCategory,
					Position: 100,
				})
				if err != nil {
					utilities.InteractionReply(s, i, "カテゴリの作成に失敗しました:"+err.Error())
					return
				}
				guild.CategoryID = newCategory.ID
				err = controller.GuildInteractor.Update(&guild)
				if err != nil {
					utilities.InteractionReply(s, i, "カテゴリ作成のDBへの同期に失敗しました:"+err.Error())
					return
				}
			}

			// チャンネルを作成する
			channel, err := s.GuildChannelCreateComplex(i.GuildID, discordgo.GuildChannelCreateData{
				Name:     channelName,
				Type:     discordgo.ChannelTypeGuildText,
				Topic:    tag,
				ParentID: guild.CategoryID,
				PermissionOverwrites: []*discordgo.PermissionOverwrite{
					{
						ID:   i.GuildID,
						Type: discordgo.PermissionOverwriteTypeRole,
						Allow: discordgo.PermissionReadMessages |
							discordgo.PermissionReadMessageHistory,
						Deny: discordgo.PermissionSendMessages,
					},
					{
						ID:   s.State.User.ID,
						Type: discordgo.PermissionOverwriteTypeMember,
						Allow: discordgo.PermissionReadMessages |
							discordgo.PermissionReadMessageHistory |
							discordgo.PermissionSendMessages,
					},
				},
			})
			if err != nil {
				utilities.InteractionReply(s, i, "チャンネルの作成に失敗しました:"+err.Error())
				return
			}

			// データの保存
			c := domain.Channel{
				ID:             channel.ID,
				Name:           channel.Name,
				GuildID:        channel.GuildID,
				BotID:          i.Interaction.AppID,
				Searchword:     tag,
				LastSearchedAt: time.Now(),
			}

			err = controller.ChannelInteractor.Create(c)
			if err != nil {
				utilities.InteractionReply(s, i, "データの保存に失敗しました。:"+err.Error())
				return
			}

			// チャンネルの作成完了通知
			utilities.InteractionReply(s, i, fmt.Sprintf("チャンネル `%s` を作成しました（タグ: `%s`）", channelName, tag))
		}

	}
}

func (controller *ChannelController) Update(s *discordgo.Session, c *discordgo.ChannelUpdate) {

	// DBにあるサーバー情報を取得
	guild, err := controller.GuildInteractor.FetchOneById(c.GuildID)

	if err != nil {
		s.ChannelMessageSend(guild.AdminChannelID, "予期せぬエラーが発生しました。:"+err.Error())
		return
	}
	// 変更されたチャンネルのIDを取得する
	updatedChannelID := c.Channel.ID
	channel, err := controller.ChannelInteractor.FetchOneById(updatedChannelID)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return
	}

	if channel.Name == c.Channel.Name && channel.Searchword == c.Channel.Topic {
		return
	}
	oldName := channel.Name
	oldSearchword := channel.Name

	channel.Name = c.Channel.Name
	channel.Searchword = c.Channel.Topic

	err = controller.ChannelInteractor.Update(&channel)

	if err != nil {
		// DBにあるサーバー情報を取得
		guild, err := controller.GuildInteractor.FetchOneById(c.GuildID)

		s.ChannelMessageSend(guild.AdminChannelID, "データの更新に失敗しました。:"+err.Error())
		return
	}

	s.ChannelMessageSend(guild.AdminChannelID, "チャンネル情報を変更しました:\n"+"チャンネル名: "+oldName+" → "+channel.Name+"\n"+"検索キーワード: "+oldSearchword+" → "+channel.Searchword)

}

func (controller *ChannelController) Delete(s *discordgo.Session, c *discordgo.ChannelDelete) {
	// 削除されたチャンネルのIDを取得する
	deletedChannelID := c.Channel.ID
	deletedChannelName := c.Channel.Name

	// DBにあるサーバー情報を取得
	guild, err := controller.GuildInteractor.FetchOneById(c.GuildID)

	if err != nil {
		s.ChannelMessageSend(guild.AdminChannelID, "データの取得に失敗しました。:"+err.Error())
		return
	}

	// admin-channelが消された場合
	if guild.AdminChannelID == deletedChannelID {
		guild.AdminChannelID = ""
		controller.GuildInteractor.Update(&guild)
		return
	}
	// カテゴリが消された場合
	if guild.CategoryID == deletedChannelID {
		guild.CategoryID = ""
		controller.GuildInteractor.Update(&guild)
		return
	}

	// 該当チャンネルがDBに存在するか確認
	_, err = controller.ChannelInteractor.FetchOneById(deletedChannelID)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return
	}

	err = controller.ChannelInteractor.Delete(deletedChannelID)

	if err != nil {
		s.ChannelMessageSend(guild.AdminChannelID, "データの更新に失敗しました。:"+err.Error())
		return
	}

	s.ChannelMessageSend(guild.AdminChannelID, "チャンネルを削除しました:\n"+"チャンネル名: "+deletedChannelName)
}
