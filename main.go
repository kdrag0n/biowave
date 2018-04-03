package main

import (
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"

	"github.com/getsentry/raven-go"
	"github.com/kdrag0n/biowave/core"
	_ "github.com/kdrag0n/biowave/modules"
)

var (
	dirInit = [...]string{
		"data",
		"data/db",
	}
)

func ensureDirectories() {
	for _, dir := range dirInit {
		ensureDirectory(dir)
	}
}

func ensureDirectory(dir string) {
	stat, err := os.Stat(dir)

	if os.IsNotExist(err) {
		err = os.Mkdir(dir, os.FileMode(0755))
		if err != nil {
			core.Log.Fatal("error creating directory", zap.String("path", dir), zap.Error(err))
		}
	} else if stat != nil && !stat.IsDir() {
		err = os.Remove(dir)
		if err != nil {
			core.Log.Fatal("error removing file", zap.String("path", dir), zap.Error(err))
		}

		err = os.Mkdir(dir, os.FileMode(0755))
		if err != nil {
			core.Log.Fatal("error creating directory", zap.String("path", dir), zap.Error(err))
		}
	} else if err != nil {
		core.Log.Error("error checking directory", zap.Error(err))
	}
}

func main() {
	core.Log.Info("starting")

	config, err := core.LoadConfig("config.yml")
	if err != nil {
		core.Log.Fatal("error loading config", zap.Error(err))
	}

	if config.Keys.Sentry != "" {
		raven.SetDSN(config.Keys.Sentry)
		config.Sentry = true
	}

	// ensure directories are present
	ensureDirectories()
	ensureDirectory(config.DatabasePath)

	client, err := core.NewClient(config)
	if err != nil {
		core.Log.Fatal("error creating client", zap.Error(err))
	}

	defer func() {
		if !client.IsDBClosed {
			err = client.DB.Close()
			if err != nil {
				core.Log.Error("error closing database", zap.Error(err))
			}
		}
	}()
	defer func() {
		core.Log.Info("stopping client")
		client.Stop()

		core.Log.Info("cleaning up modules")
		core.ModuleCleanup()
	}()

	client.Start()

	core.Log.Info("start complete")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-sc
}
