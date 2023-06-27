# local 環境構築

1. docker起動＆コンテナログイン
```
docker compose up -d
docker compose exec go bash
```

2. envファイル追加

.envをコピーして適宜修正
```
cp .env.example .env
```

3. package インストール

```
go get

```

4. migration
```
migrate -database "mysql://default:secret@tcp(mariadb:3306)/default" -path db/migrations up

```
5. seeding
```
go run db/seeds/init.go db/seeds/bot.go db/seeds/guild.go
```

6. コマンド登録
```
chmod +x command-up.sh && ./command-up.sh

<!-- 削除時 -->
chmod +x command-down.sh && ./command-down.sh
```

7. bot起動
```
go run main.go
```



※/create-channelにて、
ChannelController.UpdateとGuildController.Createが競合しておかしいので修正必要