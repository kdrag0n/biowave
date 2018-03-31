package core

import (
	"strings"
	"strconv"
	"reflect"
	"unsafe"
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

// StrID converts a snowflake ID from uint64 to string.
func StrID(id uint64) string {
	if id == 0 {
		return "@me"
	}

	return strconv.FormatUint(id, 10)
}

// ParseID parses a snowflake ID from a string into uint64.
func ParseID(id string) uint64 {
	idUint, _ := strconv.ParseUint(id, 10, 64) // #nosec
	return idUint
}

// BytesToString unsafely converts a []byte to a string. DO NOT MODIFY the []byte!
func BytesToString(b []byte) string {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := reflect.StringHeader{Data: bh.Data, Len: bh.Len}
	return *(*string)(unsafe.Pointer(&sh))
}

// StringToBytes unsafely converts a string to a []byte. DO NOT MODIFY the []byte!
func StringToBytes(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{Data: sh.Data, Len: sh.Len, Cap: sh.Len}
	return *(*[]byte)(unsafe.Pointer(&bh))
}
