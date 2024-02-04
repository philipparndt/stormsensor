package mqtt

import (
	"github.com/philipparndt/mqtt-gateway/config"
	"github.com/philipparndt/mqtt-gateway/mqtt"
)

func Start(config config.MQTTConfig) {
	mqtt.Start(config, "es_mqtt")
}
