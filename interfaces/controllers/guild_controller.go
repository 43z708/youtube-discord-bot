package controllers

import (
	"app/domain"
	"app/interfaces/database"
	"app/usecase"
	"app/utilities"
	"log"

	"github.com/bwmarrin/discordgo"
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

func (controller *GuildController) Create(s *discordgo.Session, i *discordgo.InteractionCreate) {

}

func (controller *GuildController) Update(s *discordgo.Session, i *discordgo.InteractionCreate) {

	// コマンドの場合（最初に発火）
	if i.Type == discordgo.InteractionApplicationCommand {

		// スラッシュコマンドのデータを取得する
		command := i.ApplicationCommandData()

		// /init-youtube-bot コマンド以外は無視する
		if command.Name != "init-youtube-bot" {
			return
		}

		// ユーザーが管理者であるかをチェックする
		isAdmin := utilities.IsAdmin(s, i)

		// 管理者以外のユーザーは処理を終了する
		if !isAdmin {
			utilities.InteractionReply(s, i, "このコマンドは管理者のみが実行できます。")
			return
		}

		// youtube api keyを取得
		apiKey := command.Options[0].StringValue()

		// DBにあるサーバー情報を取得
		guild, err := controller.GuildInteractor.FetchOneById(i.GuildID)

		if err != nil {
			s.ChannelMessageSend(i.ChannelID, "データの取得に失敗しました。:"+err.Error())
			return
		}
		guild.YoutubeApiKey = apiKey

		err = controller.GuildInteractor.Update(&guild)

		if err != nil {
			s.ChannelMessageSend(i.ChannelID, "データの更新に失敗しました。:"+err.Error())
			return
		}

		// カテゴリ名がすでに登録されている場合は終了
		if guild.CategoryName != "" {
			s.ChannelMessage(i.ChannelID, "Botの準備が完了しました。\"/create-channel\"コマンドでbotが投稿するチャンネルと検索ワードを設定・作成すれば設定完了です。")
			return
		}

		// カテゴリ入力を送信
		s.ChannelMessageSendComplex(i.ChannelID, &discordgo.MessageSend{
			Content: "YouTubeの動画URLを投稿するチャンネルの親カテゴリの名称を入力してください",
			Components: []discordgo.MessageComponent{
				&discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						&discordgo.TextInput{
							CustomID:    "category",
							Placeholder: "カテゴリ名をここに入力",
						},
					},
				},
			},
		})
	} else if i.Type == discordgo.InteractionMessageComponent {
		if i.MessageComponentData().CustomID == "category" {

			// ユーザーがテキスト入力を行ったときの処理をここに書く
			category := i.MessageComponentData().Values[0]

			// DBにあるサーバー情報を取得
			guild, err := controller.GuildInteractor.FetchOneById(i.GuildID)

			if err != nil {
				s.ChannelMessageSend(i.ChannelID, "データの取得に失敗しました。:"+err.Error())
				return
			}
			guild.CategoryName = category

			err = controller.GuildInteractor.Update(&guild)

			if err != nil {
				s.ChannelMessageSend(i.ChannelID, "データの更新に失敗しました。:"+err.Error())
				return
			}
			// 完了メッセージ
			s.ChannelMessage(i.ChannelID, "Botの準備が完了しました。\"/create-channel\"コマンドでbotが投稿するチャンネルと検索ワードを設定・作成すれば設定完了です。")
		}
	}
}

func (controller *GuildController) Delete(s *discordgo.Session, c *discordgo.ChannelDelete) {
}
