package utilities

import (
	"fmt"
	"net/http"
	"regexp"
)

func YoutubeValidation(Distributor string) (bool, string) {

	// 正規表現パターンを作成
	re := regexp.MustCompile(`\s+`)
	Distributor = re.ReplaceAllString(Distributor, " ")
	re = regexp.MustCompile(`^(www\.youtube\.com/@|https://www\.youtube\.com/@|youtube\.com/@|https://youtube\.com/@|@).*$`)

	// 文字列が正規表現にマッチするかチェック
	if re.MatchString(Distributor) {
		fmt.Println("Valid string")

		// 文字列を整形
		if !regexp.MustCompile(`^https://www\.`).MatchString(Distributor) {
			if regexp.MustCompile(`^@`).MatchString(Distributor) {
				Distributor = "https://www.youtube.com/" + Distributor
			} else if regexp.MustCompile(`^youtube\.com/@`).MatchString(Distributor) {
				Distributor = "https://www." + Distributor
			} else if regexp.MustCompile(`^https://youtube\.com/@`).MatchString(Distributor) {
				Distributor = "https://www." + Distributor[8:]
			}
		}

		// URLが有効かどうかチェックします。
		resp, err := http.Get(Distributor)
		if err != nil {
			// エラーがある場合はそのエラーメッセージを返す
			return false, fmt.Sprintf("当該チャンネルにアクセスできません。: %s", err)
		}
		defer resp.Body.Close()

		// HTTPステータスコードをチェックする
		if resp.StatusCode == http.StatusNotFound {
			// ステータスコードが404の場合はエラーメッセージを返す
			return false, "当該チャンネルが見つかりません。"
		} else if resp.StatusCode != http.StatusOK {
			// ステータスコードが200以外の場合はエラーメッセージを返す
			return false, fmt.Sprintf("当該チャンネルのURLが不正です。ステータスコード: %d", resp.StatusCode)
		}
		return true, Distributor

	} else {
		return false, fmt.Sprintf("正しいyoutubeチャンネルを指定してください。")
	}
}
