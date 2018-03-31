package core

import (
	"strings"
	"strconv"
)

const (
	truncateSuffix     = "..."
	codeTruncateSuffix = "```" + truncateSuffix
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

// Truncate truncates a message to 2000 characters and accounts for code blocks.
func Truncate(message string) string {
	if len(message) < 2000 {
		return message
	}

	m := message[:2000]
	if strings.Count(m, "```")%2 == 1 {
		return m[:len(m)-len(codeTruncateSuffix)] + codeTruncateSuffix
	}

	return m[:len(m)-len(truncateSuffix)] + truncateSuffix
}

func idToStr(id uint64) string {
	if id == 0 {
		return "@me"
	}

	return strconv.FormatUint(id, 10)
}

func strToID(id string) uint64 {
	idUint, _ := strconv.ParseUint(id, 10, 64)
	return idUint
}
