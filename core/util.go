package core

import (
	"strings"
)

func min(a int, b int) int {
	if a < b {
		return a
	}

	return b
}

// FilterMessage filters @everyone and @here mentions in a message to not mention.
func FilterMessage(message string) string {
	return strings.Replace(strings.Replace(message, "@everyone", "@\u200beveryone", -1), "@here", "@\u200bhere", -1)
}

