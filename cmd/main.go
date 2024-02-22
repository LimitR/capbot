package main

import (
	"capbot/internal/bot"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	b, err := bot.NewBot(os.Getenv("TOKEN"))
	if err != nil {
		panic(err)
	}
	b.Run()
}
