package wind

import (
	"encoding/json"
	"fmt"
	"github.com/philipparndt/go-logger"
	"github.com/philipparndt/mqtt-gateway/mqtt"
	"github.com/philipparndt/stormsensor/config"
	"log"
	"time"
)

var timer *time.Timer
var initialized bool
var cfg config.Config

type MessageType struct {
	Wind float64 `json:"wind"`
}

func Start(config config.Config) {
	cfg = config
	go run()
}

func onMessage(_ string, bytes []byte) {
	logger.Debug("Received wind data:", string(bytes))
	var wind MessageType
	var err = json.Unmarshal(bytes, &wind)
	if err != nil {
		logger.Error("Error unmarshalling wind data:", err, string(bytes))
		return
	}

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
		logger.Info(
			fmt.Sprintf("Detected storm, starting timer. Storm will be disabled in minutes (or later): %s",
				formatDuration(duration)))
		mqtt.PublishAbsolute(cfg.MQTT.Topic, "true", false)
	} else {
		logger.Info("Still storm, resetting timer")
		timer.Stop()
	}

	timer = time.AfterFunc(duration, func() {
		logger.Info("Timer expired, disable storm mode")
		mqtt.PublishAbsolute(cfg.MQTT.Topic, "false", false)
		timer = nil
	})
}

func consumeWind(wind MessageType) {
	windSpeed := cfg.Storm.WindSpeed

	if wind.Wind >= windSpeed {
		logger.Info(
			fmt.Sprintf("Wind speed %d exceeds threshold %d, resetting timer",
				int(wind.Wind), int(windSpeed)))

		resetTimer()
	} else if !initialized {
		initialized = true
		log.Printf("Initialized with wind speed %v\n", wind)

		logger.Info(
			fmt.Sprintf("Initialized with wind speed %d",
				int(wind.Wind)))

		mqtt.PublishAbsolute(cfg.MQTT.Topic, "false", false)
	}

	logger.Debug(fmt.Sprintf("Consumed wind data: %f", wind.Wind))
}

func run() {
	mqtt.Subscribe(cfg.Storm.WindTopic, onMessage)

	select {}
}
