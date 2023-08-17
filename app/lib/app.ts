import { log } from "./logger"
import { connectMqtt } from "./mqtt/mqtt-client"

export const startApp = async () => {
    try {
        const mqttCleanUp = await connectMqtt()
        log.info("Application is now ready.")

        return () => {
            mqttCleanUp()
        }
    }
    catch (e) {
        log.error("Application failed to start", e)
        process.exit(1)
    }
}
