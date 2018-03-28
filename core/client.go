package core

import (
	"github.com/kdrag0n/cyborg/platform"
)

// Client is connected to the platform
type Client struct {
	Platform platform.Adapter
}

// NewClient creates a new platform-independent client
func NewClient(platform platform.Adapter) *Client {
	return &Client{
		Platform: platform,
	}
}

// AddHandler registers an event handler to be invoked.
// The returned func() removes the event handler.
func (c *Client) AddHandler(handler interface{}) func() {

	return func() {
		
	}
}

// Start initiates the client's connection
func (c *Client) Start() error {
	return c.Platform.Connect()
}
