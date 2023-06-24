package main

import (
	"app/infrastructure"
)

func main() {
	db := infrastructure.Init()

	BotSeeds(db)
	GuildSeeds(db)
}
