package modules

// Initialize modules
import (
	// Essential: basic commands expected to exist
	_ "github.com/kdrag0n/biowave/modules/essential"

	// Developer: commands for developer use only
	_ "github.com/kdrag0n/biowave/modules/developer"
	
	// Images: imaging, cats, memes, profiles, etc
	_ "github.com/kdrag0n/biowave/modules/images"

	// Miscellaneous: assorted stuff
	_ "github.com/kdrag0n/biowave/modules/misc"

	// Moderation: useful tools for moderators. kick, ban, warn, etc
	_ "github.com/kdrag0n/biowave/modules/moderation"

	// Music: music in voice channels from youtube, soundcloud, etc
	_ "github.com/kdrag0n/biowave/modules/music"

	// Search: search various platforms/sources. google, wikipedia, pokemon, etc
	_ "github.com/kdrag0n/biowave/modules/search"

	// Time: time related stuff. reminders, polls, etc
	_ "github.com/kdrag0n/biowave/modules/time"
)
