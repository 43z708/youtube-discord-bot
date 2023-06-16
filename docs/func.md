# 動作

## 1 時間おき(可変)に cron で以下を実行

1. db から全チャンネルの serachword を取得する
2. youtube api search で searchword を検索
3. response が異常系の場合は admin チャンネルにエラーログを流す
4. 正常系の場合、blacklist と照合し、問題ない動画を movies の url を該当のチャンネルに投稿
5. movies に登録

※channels/searchedAt を使うことで、無駄な api を叩いたり、肥大化しやすい movies テーブルの論理削除をしやすいようにする

## 1 週間おきなどに 1 度 batch 処理

1. channels/searchedAt 以前
2. movies で直近 1 週間分を論理削除

## discord のコマンド

- ブラックリストの登録 /blacklist チャンネル URL で登録
