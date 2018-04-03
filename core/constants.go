package core

import (
	"github.com/kdrag0n/discordgo"
)

// Guild IDs
const (
	GuildPrivate uint64 = 250780048943087618
)

// User IDs
const (
	UserProduction uint64 = 239775420470394897
	UserOriginalOwner uint64 = 160567046642335746
)

// Game types
const (
	Playing discordgo.GameType = iota
	Streaming
	Listening
	Watching
)

// Lengths
const (
	LenUint64 = 8
)

// Lists
var (
	appPkgPrefixes = []string{"github.com/kdrag0n"}
)
