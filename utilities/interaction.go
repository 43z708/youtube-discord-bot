package utilities

import "github.com/bwmarrin/discordgo"

func IsAdmin(s *discordgo.Session, i *discordgo.InteractionCreate) bool {

	// ユーザーが管理者であるかをチェックする
	isAdmin := false
	for _, role := range i.Interaction.Member.Roles {
		roleInfo, err := s.State.Role(i.GuildID, role)
		if err != nil {
			// ロール情報の取得に失敗した場合はエラーハンドリングする
			// エラーメッセージを返すなどの処理を行う
			return false
		}

		//  roleInfo.Permissions のビットマスクに管理者権限が含まれているかチェック
		if roleInfo.Permissions&discordgo.PermissionAdministrator != 0 {
			isAdmin = true
			break
		}
	}

	return isAdmin
}

func InteractionReply(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
