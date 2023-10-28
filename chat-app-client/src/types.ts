export interface TMessage {
    from: string
    content: string
    id: string
    channelID: string
    timestamp: number
}

export interface TChannel {
    id: string
    members: string[]
    lastActivity: TMessage
}