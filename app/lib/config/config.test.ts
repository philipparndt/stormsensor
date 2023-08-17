import path from "path"
import { log } from "../logger"
import { applyDefaults, loadConfig } from "./config"

describe("Config", () => {
    test("default values", async () => {
        const config = {
            mqtt: {
                url: "tcp://192.168.1.1:1883",
                topic: "storm"
            }
        }
        expect(applyDefaults(config)).toStrictEqual({
            "send-full-update": true,
            loglevel: "info",
            mqtt: {
                qos: 1,
                retain: true,
                "bridge-info": true,
                url: "tcp://192.168.1.1:1883",
                topic: "storm"
            }
        })
    })

    test("load from file", () => {
        loadConfig(path.join(__dirname, "../../../production/config/config-example.json"))
        log.off()
    })
})
