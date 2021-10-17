package discord

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	dyn "MelvinBot/src/dynamo"
	"MelvinBot/src/stats"

	disc "github.com/bwmarrin/discordgo"
)

type Bot struct {
	discord *disc.Session
	dynamo  *dyn.DynamoClient
}

func NewBot(token string) Bot {
	discord, err := disc.New("Bot " + token)
	if err != nil {
		log.Fatal("could not connect to discord")
	}

	dynamo, err := dyn.NewDynamoSession()
	if err != nil {
		log.Fatal("could not connect to dynamo")
	}

	return Bot{discord, &dynamo}
}

func (bot Bot) RunBot() {

	err := bot.dynamo.GetStatsOnAllGuilds()
	if err != nil {
		log.Fatal(err)
	}
	statsChan := bot.dynamo.PutStatsOnTimer(5 * time.Minute)

	// Add handlers here
	bot.discord.AddHandler(monkaS)
	bot.discord.AddHandler(stats.TrackStats)
	bot.discord.AddHandler(stats.PrintStats)
	bot.discord.AddHandler(pinFromReaction)
	bot.discord.AddHandler(unpinFromReaction)

	err = bot.discord.Open()
	if err != nil {
		log.Fatal("couldnt open connection", err)
	}
	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// End the goroutine that is updating stats
	statsChan <- true

	// Place stats one last time for consistency
	err = bot.dynamo.PutStatsOnAllGuilds()
	if err != nil {
		log.Printf("failed dynamo put call on shutdown: %v", err)
	}
	// Cleanly close down the Discord session.
	bot.discord.Close()
}

// Handlers
func monkaS(s *disc.Session, m *disc.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return // it me
	}

	if strings.Contains(strings.ToLower(m.Message.Content), "monkas") {
		s.ChannelMessageSend(m.ChannelID, "monkaS")
	}
}
