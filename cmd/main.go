package main

import (
	"log"

	"github.com/rawly/rss-xml/pkg/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
