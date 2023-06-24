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

type ChannelController struct {
	Interactor usecase.ChannelInteractor
}

func NewChannelController(sqlHandler *gorm.DB) *ChannelController {
	return &ChannelController{
		Interactor: usecase.ChannelInteractor{
			ChannelRepository: &database.ChannelRepository{
				SqlHandler: sqlHandler,
			},
		},
	}
}

func (controller *ChannelController) FetchOneById(id string) domain.Channel {
	channel, err := controller.Interactor.FetchOneById(id)

	if err != nil {
		log.Fatalf("Error getting channels data: %s", err.Error())
	}
	return channel
}

func (controller *ChannelController) FetchAllByBotID(id string) domain.Channels {
	channels, err := controller.Interactor.FetchAllByBotID(id)
	if err != nil {
		log.Fatalf("Error getting channels data: %s", err.Error())
	}

	return channels
}

func (controller *ChannelController) FetchAllByGuildID(id string) domain.Channels {
	channels, err := controller.Interactor.FetchAllByGuildID(id)
	if err != nil {
		log.Fatalf("Error getting channels data: %s", err.Error())
	}

	return channels
}

func (controller *ChannelController) Create(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading dotenv: %s", err.Error())
	}

	// スラッシュコマンドのデータを取得する
	command := i.ApplicationCommandData()

	// /channel コマンド以外は無視する
	if command.Name != "create-channel" {
		return
	}

	// ユーザーが管理者であるかをチェックする
	isAdmin := utilities.IsAdmin(s, i)

	// 管理者以外のユーザーは処理を終了する
	if !isAdmin {
		utilities.InteractionReply(s, i, "このコマンドは管理者のみが実行できます。")
		return
	}

	// チャンネル名とタグのパラメータを取得する
	channelName := command.Options[0].StringValue()
	tag := command.Options[1].StringValue()

	// サーバー内の全チャンネルを取得
	channels, _ := s.GuildChannels(i.GuildID)

	// youtubeリンク投稿カテゴリの存在チェック
	CATEGORY_NAME := os.Getenv("CATEGORY_NAME")
	var youtubeCategoryID string
	for _, channel := range channels {
		if channel.Type == discordgo.ChannelTypeGuildCategory && channel.Name == CATEGORY_NAME {
			youtubeCategoryID = channel.ID
		}
	}
	// 存在しない場合カテゴリを作成
	if youtubeCategoryID == "" {
		newCategory, err := s.GuildChannelCreateComplex(i.GuildID, discordgo.GuildChannelCreateData{
			Name:     CATEGORY_NAME,
			Type:     discordgo.ChannelTypeGuildCategory,
			Position: 100,
		})
		if err != nil {
			utilities.InteractionReply(s, i, "カテゴリの作成に失敗しました:"+err.Error())
			return
		}
		youtubeCategoryID = newCategory.ID
	}

	// チャンネルを作成する
	channel, err := s.GuildChannelCreateComplex(i.GuildID, discordgo.GuildChannelCreateData{
		Name:     channelName,
		Type:     discordgo.ChannelTypeGuildText,
		Topic:    tag,
		ParentID: youtubeCategoryID,
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
		ID:         channel.ID,
		Name:       channel.Name,
		GuildID:    channel.GuildID,
		BotID:      i.Interaction.AppID,
		Searchword: tag,
	}

	err = controller.Interactor.Create(c)
	if err != nil {
		utilities.InteractionReply(s, i, "データの保存に失敗しました。:"+err.Error())
		return
	}

	// チャンネルの作成完了通知
	utilities.InteractionReply(s, i, fmt.Sprintf("チャンネル `%s` を作成しました（タグ: `%s`）", channelName, tag))
}

func (controller *ChannelController) Update(s *discordgo.Session, c *discordgo.ChannelUpdate) {
	// 変更されたチャンネルのIDを取得する
	updatedChannelID := c.Channel.ID
	channel, err := controller.Interactor.FetchOneById(updatedChannelID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return
		}
		s.ChannelMessageSend(channel.ID, "予期せぬエラーが発生しました。:"+err.Error())
		return
	}

	if channel.Name == c.Channel.Name && channel.Searchword == c.Channel.Topic {
		return
	}
	oldName := channel.Name
	oldSearchword := channel.Name

	channel.Name = c.Channel.Name
	channel.Searchword = c.Channel.Topic

	err = controller.Interactor.Update(channel)

	if err != nil {
		s.ChannelMessageSend(channel.ID, "データの更新に失敗しました。:"+err.Error())
		return
	}

	s.ChannelMessageSend(channel.ID, "チャンネル情報を変更しました:\n"+"チャンネル名: "+oldName+" → "+channel.Name+"\n"+"検索キーワード: "+oldSearchword+" → "+channel.Searchword)

}

func (controller *ChannelController) Delete(s *discordgo.Session, c *discordgo.ChannelDelete) {
	// 削除されたチャンネルのIDを取得する
	deletedChannelID := c.Channel.ID

	err := controller.Interactor.Delete(deletedChannelID)

	if err != nil {
		fmt.Println("Error has occured:", err)
		return
	}

	fmt.Println("Channel deleted:", deletedChannelID)
}
