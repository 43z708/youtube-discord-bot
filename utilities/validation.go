package utilities

import (
	"fmt"
	"net/http"
)

func YoutubeValidation(Distributor string) (bool, string) {

	// URLが有効かどうかチェックします。
	resp, err := http.Get("https://www.youtube.com/channel/" + Distributor)
	if err != nil {
		// エラーがある場合はそのエラーメッセージを返す
		return false, fmt.Sprintf("当該チャンネルにアクセスできません。 %s : %s", ExplainGetYoutubeChannelID, err)
	}
	defer resp.Body.Close()

	// HTTPステータスコードをチェックする
	if resp.StatusCode == http.StatusNotFound {
		// ステータスコードが404の場合はエラーメッセージを返す
		return false, "当該チャンネルが見つかりません。" + ExplainGetYoutubeChannelID
	} else if resp.StatusCode != http.StatusOK {
		// ステータスコードが200以外の場合はエラーメッセージを返す
		return false, fmt.Sprintf("当該チャンネルのURLが不正です。%s ステータスコード: %d", ExplainGetYoutubeChannelID, resp.StatusCode)
	}
	return true, Distributor

}

const ExplainGetYoutubeChannelID = "https://cdn.discordapp.com/attachments/1123945478359359569/1124615463373111326/2023-07-01_17h18_15.png"
