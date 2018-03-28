package core

import (
	"time"
)

// MessageEvent is fired upon receiving a message
type MessageEvent struct {
	Sender  *User
	Message *Message
}

// A MessageEventHandler handles MessageEvents
type MessageEventHandler func(MessageEvent)

// ReadyEvent is fired when the client is ready to act
type ReadyEvent struct {
	time time.Time
}

// A ReadyEventHandler handles ReadyEvents
type ReadyEventHandler func(ReadyEvent)
