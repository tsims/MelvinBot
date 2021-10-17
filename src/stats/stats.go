package stats

import (
	"fmt"
	"sort"
	"sync"

	disc "github.com/bwmarrin/discordgo"
)

type Stats struct {
	StatMap map[string]int
	Lock    *sync.Mutex
}

var StatsPerGuild map[string]*Stats = map[string]*Stats{}

var MelvinIDToUsernameMap sync.Map

func TrackStats(s *disc.Session, m *disc.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return // it me
	}

	guildStats, ok := StatsPerGuild[m.GuildID]
	if !ok {
		guildStats = &Stats{
			StatMap: map[string]int{},
			Lock:    &sync.Mutex{},
		}
		// What are the chances this races
		StatsPerGuild[m.GuildID] = guildStats
	}

	guildStats.Lock.Lock()
	defer guildStats.Lock.Unlock()

	user, known := MelvinIDToUsernameMap.Load(m.Author.ID)
	if !known || user != m.Author.Username {
		MelvinIDToUsernameMap.Store(m.Author.ID, m.Author.Username)
	}

	_, ok = guildStats.StatMap[m.Author.ID]
	if ok {
		guildStats.StatMap[m.Author.ID]++
	} else {
		guildStats.StatMap[m.Author.ID] = 1
	}
}

func PrintStats(s *disc.Session, m *disc.MessageCreate) {
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

	guildStats.Lock.Lock()

	// Sort and create stats array
	type MelvinPosts struct {
		name  string
		posts int
	}

	sortable := []MelvinPosts{}
	for melvinID, posts := range guildStats.StatMap {
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
	guildStats.Lock.Unlock()

	sort.Slice(sortable, func(i, j int) bool {
		return sortable[i].posts > sortable[j].posts
	})

	statsMessage := "Melvin Posts Leaderboard:"
	for _, message := range sortable {
		statsMessage += fmt.Sprintf("\n%s : %d", message.name, message.posts)
	}

	s.ChannelMessageSend(m.ChannelID, statsMessage)
}
