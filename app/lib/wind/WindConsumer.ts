import Duration from "@icholy/duration"
import { getAppConfig } from "../config/config"
import { log } from "../logger"
import { publishAbsolute } from "../mqtt/mqtt-client"

type WindType = {
    wind: number
}

const isWind = (data: any): data is WindType => {
    return data && typeof data.wind === "number"
}

export const onData = (message: Buffer) => {
    try {
        const data = JSON.parse(message.toString())
        if (isWind(data)) {
            consumeWind(data)
        }
        else {
            log.warn("Unknown message", message.toString())
        }
    }
    catch (e) {
        log.error("Failed to consume data", e, message.toString())
    }
}

let timer: (ReturnType<typeof setTimeout> | undefined)
let initialized = false

const resetTimer = () => {
    const config = getAppConfig()
    const duration = Duration.seconds(config.storm.resetTimeSeconds)

    if (!timer) {
        log.info("Detected storm, starting timer. Storm will be disabled in minutes (or later):", duration.minutes)
        publishAbsolute(true, config.mqtt.topic)
    }
    else {
        log.debug("Still storm, resetting timer")
        clearTimeout(timer)
    }

    timer = setTimeout(() => {
        log.info("Timer expired, disable storm mode")
        publishAbsolute(false, config.mqtt.topic)
        timer = undefined
    }, duration.milliseconds())
}

const consumeWind = (wind: WindType) => {
    const storm = getAppConfig().storm
    if (wind.wind >= storm.windSpeed) {
        log.debug(`Wind speed ${wind.wind} exceeds threshold ${storm.windSpeed}, resetting timer`)
        resetTimer()
    }
    else if (!initialized) {
        initialized = true
        log.info(`Initialized with wind speed ${wind.wind}`)
        const config = getAppConfig()
        publishAbsolute(false, config.mqtt.topic)
    }
}
