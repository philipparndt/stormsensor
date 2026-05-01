package wind

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/philipparndt/go-logger"
	"github.com/philipparndt/mqtt-gateway/mqtt"
	"github.com/philipparndt/stormsensor/config"
)

var timer *time.Timer
var initialized bool
var cfg config.Config
var lastWind float64
var stormActive bool

type Status struct {
	Wind        float64 `json:"wind"`
	StormActive bool    `json:"stormActive"`
	Initialized bool    `json:"initialized"`
}

func GetStatus() Status {
	return Status{
		Wind:        lastWind,
		StormActive: stormActive,
		Initialized: initialized,
	}
}

type MessageType struct {
	Wind float64 `json:"wind"`
}

func Start(config config.Config) {
	cfg = config
	go run()
}

func onMessage(_ string, bytes []byte) {
	logger.Debug("Received wind data", "data", string(bytes))
	var wind MessageType
	var err = json.Unmarshal(bytes, &wind)
	if err != nil {
		logger.Error("Error unmarshalling wind data", "error", err, "data", string(bytes))
		return
	}

	lastWind = wind.Wind
	consumeWind(wind)
}

func formatDuration(duration time.Duration) string {
	if duration.Seconds() < 60 {
		return fmt.Sprintf("%d seconds", int(duration.Seconds()))
	}

	return fmt.Sprintf("%d minutes", int(duration.Minutes()))
}

func resetTimer() {
	duration := time.Duration(cfg.Storm.ResetTimeSeconds) * time.Second

	if timer == nil {
		logger.Info("Detected storm, starting timer",
			"resetIn", formatDuration(duration))
		stormActive = true
		mqtt.PublishAbsolute(cfg.MQTT.Topic, "true", false)
	} else {
		logger.Info("Still storm, resetting timer")
		timer.Stop()
	}

	timer = time.AfterFunc(duration, func() {
		logger.Info("Timer expired, disable storm mode")
		stormActive = false
		mqtt.PublishAbsolute(cfg.MQTT.Topic, "false", false)
		timer = nil
	})
}

func consumeWind(wind MessageType) {
	windSpeed := cfg.Storm.WindSpeed

	if wind.Wind >= windSpeed {
		logger.Info("Wind speed exceeds threshold, resetting timer",
			"wind", int(wind.Wind), "threshold", int(windSpeed))

		resetTimer()
	} else if !initialized {
		initialized = true

		logger.Info("Initialized with wind speed",
			"wind", int(wind.Wind))

		mqtt.PublishAbsolute(cfg.MQTT.Topic, "false", false)
	}

	logger.Debug("Consumed wind data", "wind", wind.Wind)
}

func run() {
	mqtt.Subscribe(cfg.Storm.WindTopic, onMessage)

	select {}
}
