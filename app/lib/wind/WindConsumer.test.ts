import Duration from "@icholy/duration"
import { Writable } from "stream"
import winston from "winston"
import { applyConfig } from "../config/config"
import { createLogger, log, setLogger } from "../logger"
import { onData } from "./WindConsumer"
import * as mqtt from "../mqtt/mqtt-client"

type Message = {
    message: any
}

let messages: Message[]
jest.spyOn(mqtt, "publishAbsolute").mockImplementation((message) => {
    messages.push({ message })
})

export class TestLogger {
    output = ""
    logger: winston.Logger

    constructor () {
        process.env.FORCE_COLOR = "0"

        const stream = new Writable()
        stream._write = (chunk, _, next) => {
            this.output = this.output += chunk.toString()
            next()
        }
        this.logger = createLogger(new winston.transports.Stream({ stream }))
        this.logger.level = "TRACE"
        setLogger(this.logger)
    }
}

describe("WindConsumer", () => {
    beforeAll(() => {
        applyConfig({
            storm: {
                windSpeed: 15,
                resetTimeSeconds: 60,
                windTopic: "wind"
            }
        })
    })

    beforeEach(() => {
        messages = []
        jest.useFakeTimers()
        log.off()
    })

    afterEach(() => {
        jest.useRealTimers()
        log.on()
    })

    describe("invalid messages", () => {
        let logger: TestLogger

        beforeEach(() => {
            logger = new TestLogger()
        })

        test("no JSON", () => {
            onData(Buffer.from("hello"))
            expect(logger.output).toMatch(/\d+-\d+-\d+T\d+:\d+:\d+.\d+.* \[.*ERROR.*] Failed to consume data.*/)
        })

        test("no wind data", () => {
            onData(Buffer.from(JSON.stringify({ wind: "hello" })))
            expect(logger.output).toMatch(/\d+-\d+-\d+T\d+:\d+:\d+.\d+.* \[.*WARN.*] Unknown message.*/)
        })
    })

    test("no storm", () => {
        onData(Buffer.from(JSON.stringify({ wind: 14 })))
        expect(messages.length).toBe(0)
    })

    test("storm", () => {
        onData(Buffer.from(JSON.stringify({ wind: 31 })))
        expect(messages.length).toBe(1)
        expect(messages[0].message).toBe(true)
        jest.advanceTimersByTime(Duration.seconds(65).milliseconds())
        expect(messages.length).toBe(2)
        expect(messages[1].message).toBe(false)
    })

    test("keep storm", () => {
        for (let i = 0; i < 10; i++) {
            onData(Buffer.from(JSON.stringify({ wind: 31 })))
            expect(messages.length).toBe(1)
            expect(messages[0].message).toBe(true)
            jest.advanceTimersByTime(Duration.seconds(30).milliseconds())
        }

        jest.advanceTimersByTime(Duration.seconds(61).milliseconds())
        expect(messages.length).toBe(2)
        expect(messages[1].message).toBe(false)
    })
})
