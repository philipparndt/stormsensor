package main

import (
	"expvar"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/philipparndt/go-logger"
	"github.com/philipparndt/stormsensor/config"
	"github.com/philipparndt/stormsensor/mqtt"
	"github.com/philipparndt/stormsensor/version"
	"github.com/philipparndt/stormsensor/wind"
)

func initPprof() {
	go func() {
		http.ListenAndServe(":6060", nil)
	}()
}

func main() {
	logger.Init("info", logger.Logger())
	logger.Info("stormsensor", "version", version.Info())

	if len(os.Args) < 2 {
		logger.Error("No config file specified")
		os.Exit(1)
	}

	configFile := os.Args[1]
	logger.Info("Config file", "file", configFile)
	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		logger.Error("Failed loading config", "error", err)
		return
	}

	logger.SetLevel(cfg.LogLevel)
	initPprof()

	expvar.Publish("stormStatus", expvar.Func(func() any {
		return wind.GetStatus()
	}))

	mqtt.Start(cfg.MQTT)
	wind.Start(cfg)

	logger.Info("Application is now ready. Press Ctrl+C to quit.")

	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel

	logger.Info("Received quit signal")
}
