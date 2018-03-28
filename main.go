package main

import (
	"github.com/kdrag0n/biowave/core"
	"github.com/getsentry/raven-go"
)

func main() {
	config, err := core.LoadConfig("config.json")
	if err != nil {
		panic(err)
	}

	if config.Keys.Sentry != "" {
		raven.SetDSN(config.Keys.Sentry)
		config.Sentry = true
	}

	raven.CapturePanic()

	client := core.NewClient(config)
	client.Start()
}
