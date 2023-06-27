# 動作
## botについて
- 1botにつき、原則1サーバー（内部的には複数サーバーに対応できるようにしておく）
- botを招待した際、DBにサーバーの登録がなければ登録する(YoutubeApiKeyとCategoryNameはnullのまま)
- botをキックした場合はDBはそのまま

## 1 時間おき(可変)に cron で以下を実行

1. db から全チャンネルの serachword を取得する
2. youtube api search で searchword を検索
3. response が異常系の場合はチャンネルにエラーログを流す
4. 正常系の場合、blacklist と照合し、問題ない動画のurl を該当のチャンネルに投稿

※channels/LastSearchedAt を使うことで、無駄な api を叩かないようにする

## 1 週間おきなどに 1 度 batch 処理

1. channels/LastSearchedAt 以前
2. movies で直近 1 週間分を論理削除

## discord のコマンド
- "/register-apikey YOUTUBEのAPIKEY" でapikeyの登録・更新
- "/create-channel チャンネル名 検索ワード" でチャンネルの作成および検索ワードの設定、変更・削除はhandlerで同期させる
※検索ワードは、チャンネルのTopic欄の内容とする。
- ブラックリストの登録 "/add-blacklist チャンネルURL" で登録
- ブラックリストの一覧 "/get-blacklist" で一覧を出力
- ブラックリストの解除 "/remove-blacklist チャンネルURL" で解除

## 管理者マニュアル
1. 導入したいサーバーにbotを招待
2. "/register-apikey YOUTUBEのAPIKEY" でapikeyの登録・更新
3. "/create-channel チャンネル名 検索ワード" でチャンネルの作成および検索ワードの設定