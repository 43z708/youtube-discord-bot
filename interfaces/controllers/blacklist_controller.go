package controllers

import (
	"app/domain"
	"app/interfaces/database"
	"app/usecase"
	"app/utilities"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

type BlacklistController struct {
	BlacklistInteractor usecase.BlacklistInteractor
	GuildInteractor     usecase.GuildInteractor
}

func NewBlacklistController(sqlHandler *gorm.DB) *BlacklistController {
	return &BlacklistController{
		BlacklistInteractor: usecase.BlacklistInteractor{
			BlacklistRepository: &database.BlacklistRepository{
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

func (controller *BlacklistController) FetchOneById(id string) domain.Blacklist {
	blacklist, err := controller.BlacklistInteractor.FetchOneById(id)

	if err != nil {
		log.Fatalf("Error getting blacklists data: %s", err.Error())
	}
	return blacklist
}

func (controller *BlacklistController) FetchAllByGuildID(id string) domain.Blacklists {
	blacklists, err := controller.BlacklistInteractor.FetchAllByGuildID(id)
	if err != nil {
		log.Fatalf("Error getting blacklists data: %s", err.Error())
	}

	return blacklists
}

func (controller *BlacklistController) FetchBlacklist(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionApplicationCommand {
		// スラッシュコマンドのデータを取得する
		command := i.ApplicationCommandData()

		// /get-blacklist コマンド以外は無視する
		if command.Name == "get-blacklist" {

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

			blacklists, err := controller.BlacklistInteractor.FetchAllByGuildID(i.GuildID)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					utilities.InteractionReply(s, i, "ブラックリストに登録がありません。\n /add-blacklistコマンドで登録が可能です。")
					return
				}
				utilities.InteractionReply(s, i, "DBからのサーバー情報の取得に失敗しました。:"+err.Error())
				return
			}

			list := make([]string, 0)
			for _, blacklist := range blacklists {
				list = append(list, blacklist.Distributor)
			}
			if len(list) > 0 {
				utilities.InteractionReply(s, i, "ブラックリストは以下です。\n"+strings.Join(list, "\n"))
			} else {
				utilities.InteractionReply(s, i, "ブラックリストに登録がありません。\n /add-blacklistコマンドで登録が可能です。")
			}
		}
	}
}

func (controller *BlacklistController) Create(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionApplicationCommand {
		// スラッシュコマンドのデータを取得する
		command := i.ApplicationCommandData()

		// /add-blacklist コマンド以外は無視する
		if command.Name == "add-blacklist" {

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
			// リクエストされたチャンネルが不正かどうかチェック
			isValidated, text := utilities.YoutubeValidation(command.Options[0].StringValue())
			if !isValidated {
				utilities.InteractionReply(s, i, text)
				return
			}

			Distributor := text

			// 重複チェック
			_, err = controller.BlacklistInteractor.FetchOneByDistributor(Distributor)
			// errがnilである場合は重複データが存在するということ
			if err == nil {
				utilities.InteractionReply(s, i, Distributor+" はすでにブラックリストに登録されています。")
				return
			}
			b := domain.Blacklist{
				Distributor: Distributor,
				BotID:       s.State.User.ID,
				GuildID:     i.GuildID,
			}
			err = controller.BlacklistInteractor.Create(b)
			if err != nil {
				utilities.InteractionReply(s, i, "データの保存に失敗しました。:"+err.Error())
				return
			}
			// 追加完了通知
			utilities.InteractionReply(s, i, fmt.Sprintf("ブラックリストに `%s` を追加しました。", Distributor))
		}
	}
}

func (controller *BlacklistController) Update(b *domain.Blacklist) {
	return
}

func (controller *BlacklistController) Delete(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionApplicationCommand {
		// スラッシュコマンドのデータを取得する
		command := i.ApplicationCommandData()

		// /remove-blacklist コマンド以外は無視する
		if command.Name == "remove-blacklist" {

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
			blacklist, err := controller.BlacklistInteractor.FetchOneByDistributor(command.Options[0].StringValue())
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					utilities.InteractionReply(s, i, command.Options[0].StringValue()+"はブラックリストに登録されていません。")
					return
				} else {
					utilities.InteractionReply(s, i, "データの取得でエラーが発生しました。")
				}
			}

			controller.BlacklistInteractor.Delete(blacklist.ID)
			utilities.InteractionReply(s, i, command.Options[0].StringValue()+" をブラックリストから削除しました。")
		}
		return
	}
}
