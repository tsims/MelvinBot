package discord

import (
	"fmt"
	"sort"
	"sync"

	disc "github.com/bwmarrin/discordgo"
)

type Stats struct {
	stats map[string]int
	lock  *sync.Mutex
}

var StatsPerGuild map[string]*Stats = map[string]*Stats{}

var MelvinIDToUsernameMap sync.Map

func trackStats(s *disc.Session, m *disc.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return // it me
	}

	guildStats, ok := StatsPerGuild[m.GuildID]
	if !ok {
		guildStats = &Stats{
			stats: map[string]int{},
			lock:  &sync.Mutex{},
		}
		// What are the chances this races
		StatsPerGuild[m.GuildID] = guildStats
	}

	guildStats.lock.Lock()
	defer guildStats.lock.Unlock()

	user, known := MelvinIDToUsernameMap.Load(m.Author.ID)
	if !known || user != m.Author.Username {
		MelvinIDToUsernameMap.Store(m.Author.ID, m.Author.Username)
	}

	_, ok = guildStats.stats[m.Author.ID]
	if ok {
		guildStats.stats[m.Author.ID]++
	} else {
		guildStats.stats[m.Author.ID] = 1
	}
}

func printStats(s *disc.Session, m *disc.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return // it me
	}

	if m.Content != "!stats" {
		return
	}

	guildStats, ok := StatsPerGuild[m.GuildID]
	if !ok {
		s.ChannelMessageSend(m.ChannelID, "Sorry I'm not tracking stats for this server")
		return
	}

	guildStats.lock.Lock()

	// Sort and create stats array
	type MelvinPosts struct {
		name  string
		posts int
	}

	sortable := []MelvinPosts{}
	for melvinID, posts := range guildStats.stats {
		username, ok := MelvinIDToUsernameMap.Load(melvinID)
		if !ok {
			continue
		}
		// haha type assertion go brr
		userString, ok := username.(string)
		if !ok {
			continue
		}
		sortable = append(sortable, MelvinPosts{name: userString, posts: posts})
	}
	// Don't need lock anymore
	guildStats.lock.Unlock()

	sort.Slice(sortable, func(i, j int) bool {
		return sortable[i].posts > sortable[j].posts
	})

	statsMessage := "Melvin Posts Leaderboard:"
	for _, message := range sortable {
		statsMessage += fmt.Sprintf("\n%s : %d", message.name, message.posts)
	}

	s.ChannelMessageSend(m.ChannelID, statsMessage)
}
