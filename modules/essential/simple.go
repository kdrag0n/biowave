package essential

import (
	"fmt"
	"github.com/kdrag0n/biowave/core"
	"time"
)

// Ping returns the time it takes to respond to a message.
func (C) Ping(c *core.Context) {
	c.Info("pong")

	before := time.Now()
	c.Loading("Pong!")

	c.EditOkFmt("Pong! %.2fms", core.Milliseconds(time.Since(before)))
}

// Uptime returns how long the bot has been running.
func (C) Uptime(c *core.Context) {
	c.Info("of the bot")

	uptime := time.Since(c.Client.StartTime)
	h := uint16(uptime.Hours())
	d := uint16(h / 24)
	m := uint16(uptime.Minutes())
	if m == 0 {
		m = 1
	}

	result := fmt.Sprintf("%dm", m)
	if h > 0 {
		result = fmt.Sprintf("%dh%s", h, result)
	}

	if d > 0 {
		result = fmt.Sprintf("%dd%s", d, result)
	}

	c.Send("I've been up for " + result + ".")
}
