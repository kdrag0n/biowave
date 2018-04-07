package time

import (
	"github.com/kdrag0n/biowave/core"
)

// RemindMe schedules reminders for users.
func (C) RemindMe(c *core.Context) {
	c.Info("schedules reminders for you")
}
