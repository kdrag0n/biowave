package essential

import (
	"fmt"
	"github.com/kdrag0n/biowave/core"
)

// Help lists commands and what they do.
func (C) Help(c *core.Context) {
	c.Info("lists commands")

	dev, _ := c.Client.Developers.GetBit(c.Event.Author.ID)
	fields := make(map[string]core.Paginator, 10)

	if true { // no args
		for _, cmd := range c.Client.Commands {
			if !(cmd.Hidden || cmd.Requires(core.PermDeveloper)) || dev {
				entry := fmt.Sprintf(" [%s](http://) %s", cmd.Name, cmd.Description)

				if pager, ok := fields[cmd.Module.Name]; ok {
					pager.AddLine(entry)
				} else {
					pager := core.NewPaginator(1024)
					pager.AddLine("\u200b" + entry)
					fields[cmd.Module.Name] = pager
				}
			}
		}
	} else {

	}

	embed := c.Embed().AuthorHead("Bot Help").RoleColor()
	for module, pages := range fields {
		for _, page := range pages.Pages() {
			embed.AddField(module, page)
		}
	}

	c.DirectSendEmbed(embed)
	c.React("✅")
}
