# 仕様書

## DB および env ファイル

### bots

| Collumn      | Type    | Meaning             |
| ------------ | ------- | ------------------- |
| id           | string  | bot の id           |
| name         | string  | bot の名前          |
| token        | string  | bot の api トークン |
| isAvalilable | boolean | 有効な bot かどうか |

※ サーバーごとに bot を作成し db に登録すること

### guilds

| Collumn | Type   | Meaning       |
| ------- | ------ | ------------- |
| id      | string | サーバーの id |
| name    | string | サーバー名    |
| botId   | string | bot の id     |

### channels

| Collumn    | Type     | Meaning                            |
| ---------- | -------- | ---------------------------------- |
| id         | string   | チャンネル ID                      |
| name       | string   | チャンネル名                       |
| guildId    | string   | サーバー ID                        |
| botId      | string   | bot の ID                          |
| searchword | string   | youtube で検索する際のハッシュタグ |
| searchedAt | string   | 最後の検索 api を叩いた日時        |
| createdAt  | datetime | 作成日時 　　　　　　　　　　　　  |
| updatedAt  | datetime | 更新日時 　　　　　　　　　　　　  |
| deletedAt  | datetime | 削除日時 　　　　　　　　　　　　  |

### movies

| Collumn     | Type     | Meaning                 |
| ----------- | -------- | ----------------------- |
| id          | string   | 動画 id                 |
| url         | string   | 動画 url                |
| distributor | string   | youtube チャンネルの id |
| botId       | string   | bot の ID               |
| guildId     | string   | サーバー ID             |
| createdAt   | datetime | 作成日時 　　　　　　　 |
| updatedAt   | datetime | 更新日時 　　　　　　　 |
| deletedAt   | datetime | 削除日時 　　　　　　　 |

### blacklists

| Collumn     | Type   | Meaning                 |
| ----------- | ------ | ----------------------- |
| id          | number | primaryKey              |
| distributor | string | youtube チャンネルの id |
| botId       | string | bot の ID               |
| guildId     | string | サーバー ID             |

### env ファイル

何時間おきか
youtube api


## コマンド等

### migrationファイルの作り方
```
migrate create -ext sql -dir db/migrations -seq ファイル名
```