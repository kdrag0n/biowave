package essential

import (
	"github.com/kdrag0n/biowave/core"
)

func init() {
	core.RegisterModule("Essential", func(m *core.Module) {
		m.Add("test", "Test the bot.", nil, cmdTest)
	})
}

func cmdTest(c *core.Context) {
	c.Send("It works!")
}
