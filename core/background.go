package core

import (
	"github.com/dgraph-io/badger"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"time"

	"github.com/kdrag0n/discordgo"
)

var printer = message.NewPrinter(language.English)

// Status represents a pair of game type and status text.
type Status struct {
	Type discordgo.GameType
	Text string
	URL  string
}

func (c *Client) housekeeper() {
	twoMin := time.Tick(2 * time.Minute)
	eightHr := time.Tick(8 * time.Hour)

	for {
		select {
		case <-twoMin:
			go c.hkTask("update status", c.UpdateStatus)
		case <-eightHr:
			go c.hkTask("clean up database", c.dbCleanup)
		}
	}
}

func (c *Client) hkTask(name string, exec func()) {
	defer c.ErrorHandler("housekeeper task: " + name)
	exec()
}

func (c *Client) dbCleanup() {
	if err := c.DB.PurgeOlderVersions(); err != nil {
		panic(err)
	}

	if err := c.DB.RunValueLogGC(0.5); err != nil && err != badger.ErrNoRewrite {
		panic(err)
	}
}

// UpdateStatus updates the status on all shards.
func (c *Client) UpdateStatus() {
	format := printer.Sprintf
	status := [...]Status{
		Status{Playing, format("with %d users", c.UserCount()), ""},
		Status{Playing, format("in %d channels", c.ChannelCount()), ""},
		Status{Playing, format("in %d servers", c.GuildCount()), ""},
		Status{Playing, "with bits and bytes", ""},
		Status{Playing, "World Domination", ""},
		Status{Playing, "with you", ""},
		Status{Playing, "with potatoes", ""},
		Status{Playing, "something", ""},
		Status{Streaming, "data", "http://hackertyper.com/"},
		Status{Streaming, "music", "https://www.youtube.com/channel/UC-9-kyTW8ZkZNDHQJ6FgpwQ"},
		Status{Streaming, "your tunes", "https://www.youtube.com/watch?v=zQJh0MWvccs"},
		Status{Listening, "you", ""},
		Status{Watching, "darkness", ""},
		Status{Watching, "streams", ""},
		Status{Streaming, "your face", "https://www.youtube.com/watch?v=IUjZtoCrpyA"},
		Status{Listening, "alone", ""},
		Status{Streaming, "Alone", "https://www.youtube.com/watch?v=YnwsMEabmSo"},
		Status{Streaming, "bits and bytes", "https://www.youtube.com/watch?v=N3ZMvqISfvY"},
		Status{Listening, "Rick Astley", ""},
		Status{Streaming, "only the very best", "https://www.youtube.com/watch?v=dQw4w9WgXcQ"},
		Status{Listening, "those potatoes", ""},
		Status{Playing, "with my fellow shards", ""},
		Status{Listening, "the cries of my shards", ""},
		Status{Streaming, "Monstercat", "https://www.twitch.tv/monstercat"},
		Status{Watching, "dem videos", ""},
		Status{Watching, "you in your sleep", ""},
		Status{Watching, "over you as I sleep", ""},
		Status{Watching, "the movement of electrons", ""},
		Status{Playing, "with some protons", ""},
		Status{Listening, "the poor electrons", ""},
		Status{Listening, "the poor neutrons", ""},
		Status{Listening, "trigger-happy players", ""},
		Status{Playing, "Discord Hacker v39.2", ""},
		Status{Playing, "Discord Hacker v42.0", ""},
		Status{Streaming, "donations", "https://patreon.com/kdragon"},
		Status{Streaming, "You should totally donate!", "https://patreon.com/kdragon"},
		Status{Listening, "my people", ""},
		Status{Listening, "my favorites", ""},
		Status{Watching, "my minions", ""},
		Status{Watching, "the chosen ones", ""},
		Status{Watching, "stars combust", ""},
		Status{Watching, "your demise", ""},
		Status{Streaming, "the supernova", "https://www.youtube.com/watch?v=5WXyCJ1w3Ks"},
		Status{Listening, "something", ""},
		Status{Streaming, "something", "https://www.youtube.com/watch?v=FM7MFYoylVs"},
		Status{Watching, "I am Cow", ""},
		Status{Watching, "you play", ""},
		Status{Watching, "for raids", ""},
		Status{Playing, "buffing before the raid", ""},
		Status{Streaming, "this sick action", "https://www.youtube.com/watch?v=tD6KJ7QtQH8"},
		Status{Listening, "memes", ""},
		Status{Watching, "memes", ""},
		Status{Watching, "that dank vid", ""},
	}

	picked := status[Rand(len(status))]

	data := discordgo.UpdateStatusData{
		Status: "online",
		AFK:    false,
		Game: &discordgo.Game{
			Name: picked.Text,
			Type: picked.Type,
			URL:  picked.URL,
		},
	}

	c.ForSessions(func(session *discordgo.Session) {
		session.UpdateStatusComplex(data)
	})
}
