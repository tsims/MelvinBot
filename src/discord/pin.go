package discord

import (
	"fmt"
	"log"

	disc "github.com/bwmarrin/discordgo"
)

// Leverage admin priveleges of the bot to look for reactions and Pin things

func pinFromReaction(s *disc.Session, m *disc.MessageReactionAdd) {
	if m.MessageReaction.Emoji.Name != "📌" {
		return
	}

	err := s.ChannelMessagePin(m.ChannelID, m.MessageID)
	if err != nil {
		log.Printf("error pinning: %v", err)
	}
}

func unpinFromReaction(s *disc.Session, m *disc.MessageReactionRemove) {
	if m.MessageReaction.Emoji.Name != "📌" {
		return
	}

	msg, err := s.ChannelMessage(m.ChannelID, m.MessageID)
	if err != nil {
		return
	}
	for _, reaction := range msg.Reactions {
		if reaction.Emoji.Name == "📌" {
			return
		}
	}

	// Check if its pinned in the first place
	var found bool
	msgs, err := s.ChannelMessagesPinned(m.ChannelID)
	if err != nil {
		log.Printf("error checking pinned msgs: %v", err)
	}
	for _, pin := range msgs {
		if pin.ID == msg.ID {
			found = true
			break
		}
	}

	if !found {
		return
	}

	err = s.ChannelMessageUnpin(m.ChannelID, m.MessageID)
	if err != nil {
		log.Printf("error unpinning: %v", err)
	}
	if err == nil {
		_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Unpinning post: %s: %s", msg.Author.Username, msg.Content))
		if err != nil {
			log.Printf("error sending unpin message: %v", err)
		}
	}
}
