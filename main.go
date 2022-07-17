package main

import (
	"os"

	"MelvinBot/src/discord"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load("/home/nelson/apps/.env")
	token := os.Getenv("token")
	bot := discord.NewBot(token)
	bot.RunBot()
}
