package main

import (
	"capbot/internal/bot"
	"capbot/internal/config"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	conf, err := config.NewConfig("./config.yaml")
	if err != nil {
		conf = config.DefaultConfig()
	}
	err = godotenv.Load()
	if err != nil {
		panic(err)
	}

	b, err := bot.NewBot(os.Getenv("TOKEN"), conf)
	if err != nil {
		panic(err)
	}
	b.Run()
}
