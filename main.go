package main

import (
	"os"

	"MelvinBot/src/discord"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	token := os.Getenv("token")
	bot := discord.NewDiscordSession(token)

	bot.RunBot()

}
