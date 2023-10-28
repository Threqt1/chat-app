import { EventEmitter } from "events"
import { TMessage } from "./types"

const WEBSOCKET_URL = "ws://localhost:8081/ws"

export default class SocketConnection {
    private socket: WebSocket
    private emitter: SocketEventEmitter
    private username: string
    private currentChannel: string

    constructor() {
        this.socket = new WebSocket(WEBSOCKET_URL)
        this.emitter = new SocketEventEmitter()
        this.username = ""
        this.currentChannel = ""
    }

    public getEmitter() {
        return this.emitter
    }

    public getUsername() {
        return this.username
    }

    public setUsername(username: string) {
        this.username = username
    }

    public getChannel() {
        return this.currentChannel
    }

    public setChannel(channel: string) {
        this.currentChannel = channel
    }

    public initialize() {
        this.socket.onmessage = message => {
            const wsMessage: WEBSOCKET_RESPONSES = JSON.parse(message.data)
            switch (wsMessage.type) {
                case "bootup":
                    this.emitter.emit("bootup", wsMessage.username)
                    break;
                case "chat":
                    if (wsMessage.message.channelID === this.currentChannel) {
                        this.emitter.emit("chat_message", wsMessage.message)
                    } else {
                        this.emitter.emit("chat_notification", wsMessage.message)
                    }
                    break
                case "channel_added":
                    this.emitter.emit("channel_added", wsMessage.channelID)
            }
        }
    }

    public connect(username: string, maxTime: number) {
        return Promise.race([new Promise((resolve, reject) => {
            this.socket.onerror = error => reject(error)

            this.emitter.once("bootup", (registeredUsername) => registeredUsername === username ? resolve(username) : reject())
            this.username = username
            this.sendBootup(username)
        }), new Promise((_, reject) => setTimeout(() => reject(), maxTime))])
    }

    private sendMessage(message: WEBSOCKET_REQUEST) {
        this.socket.send(JSON.stringify(message))
    }

    public sendBootup(username: string) {
        const bootupMessage: BootupRequest = {
            type: "bootup",
            username
        }

        this.sendMessage(bootupMessage)
    }

    public sendChat(from: string, content: string, channel: string) {
        const chatMessage: ChatRequest = {
            type: "chat",
            message: {
                from,
                content,
                channelID: channel
            },
        }

        this.sendMessage(chatMessage)
    }

    public sendChannelAdded(username: string, channelID: string) {
        const channelAddedMessage: ChannelAddedRequest = {
            type: "channel_added",
            channelID,
            username
        }

        this.sendMessage(channelAddedMessage)
    }
}

interface BaseWebsocket {
    type: string
}

interface BootupRequest extends BaseWebsocket {
    type: "bootup"
    username: string
}

interface ChatRequest extends BaseWebsocket {
    type: "chat"
    message: Omit<TMessage, "id" | "timestamp">
}

interface ChannelAddedRequest extends BaseWebsocket {
    type: "channel_added"
    username: string
    channelID: string
}

type WEBSOCKET_REQUEST = BootupRequest | ChatRequest | ChannelAddedRequest

interface BootupResponse extends BaseWebsocket {
    type: "bootup",
    username: string
}

interface ChatResponse extends BaseWebsocket {
    type: "chat",
    message: TMessage
}

interface ChannelAddedResponse extends BaseWebsocket {
    type: "channel_added"
    channelID: string
}

type WEBSOCKET_RESPONSES = BootupResponse | ChatResponse | ChannelAddedResponse

type WEBSOCKET_EMITTER_EVENTS = {
    "bootup": [username: string]
    "chat_message": [chat: TMessage]
    "chat_notification": [chat: TMessage]
    "channel_added": [channelID: string]
}

class SocketEventEmitter extends EventEmitter {
    constructor() {
        super()
    }

    public override emit<K extends keyof WEBSOCKET_EMITTER_EVENTS>(
        eventName: K,
        ...args: WEBSOCKET_EMITTER_EVENTS[K]
    ): boolean {
        return super.emit(eventName, ...args)
    }

    public override on<K extends keyof WEBSOCKET_EMITTER_EVENTS>(
        eventName: K,
        listener: (...args: WEBSOCKET_EMITTER_EVENTS[K]) => void
    ): this {
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        return super.on(eventName, listener as (...args: any[]) => void)
    }

    public override off<K extends keyof WEBSOCKET_EMITTER_EVENTS>(
        eventName: K,
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        listener: (...args: any[]) => void
    ): this {
        return super.off(eventName, listener)
    }

    public override once<K extends keyof WEBSOCKET_EMITTER_EVENTS>(
        eventName: K,
        listener: (...args: WEBSOCKET_EMITTER_EVENTS[K]) => void
    ): this {
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        return super.once(eventName, listener as (...args: any[]) => void)
    }
}