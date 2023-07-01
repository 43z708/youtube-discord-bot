package main

import (
	"app/infrastructure"
)

func main() {

	Init := infrastructure.Init()
	infrastructure.Router(Init)
	infrastructure.Discord(Init)

}
