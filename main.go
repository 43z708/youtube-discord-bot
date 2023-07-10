package main

import (
	"app/infrastructure"
)

func main() {

	Init := infrastructure.Init()
	infrastructure.Discord(Init)
	// infrastructure.Router(Init)

}
