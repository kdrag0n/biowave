package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/getsentry/raven-go"
	"github.com/kdrag0n/biowave/core"
	_ "github.com/kdrag0n/biowave/modules"
)

func main() {
	core.Log.Info("Starting.")

	config, err := core.LoadConfig("config.yml")
	if err != nil {
		panic(err)
	}

	if config.Keys.Sentry != "" {
		raven.SetDSN(config.Keys.Sentry)
		config.Sentry = true
	}

	client := core.NewClient(config)
	client.Start()

	core.Log.Info("Started.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	core.Log.Info("Stopping.")
	client.Stop()
}
