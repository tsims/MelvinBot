package discord

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	disc "github.com/bwmarrin/discordgo"
)

type Bot struct {
	discord *disc.Session
}

func NewDiscordSession(token string) Bot {
	discord, err := disc.New("Bot " + token)
	if err != nil {
		log.Fatal("could not connect to discord")
	}

	return Bot{discord}
}

func (bot Bot) RunBot() {

	// Add handlers
	bot.discord.AddHandler(monkaS)

	err := bot.discord.Open()
	if err != nil {
		log.Fatal("couldnt open connection", err)
	}
	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

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
