package main

import (
	"github.com/asb1302/innopolis_go_hw11/internal/app"
	"log"
)

func main() {
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
