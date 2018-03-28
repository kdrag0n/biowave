package platform

// Adapter provides an interface to a chat platform
type Adapter interface {
	Connect() error
	AddHandler()
}

// Config configures an adapter
type Config struct {
	Token  string `json:"token"`
	Shards int    `json:"shards"`
}
