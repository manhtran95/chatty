// WebSocketTypes.ts

// Message type definitions
export type MessageType =
    | 'ClientSendMessage'
    | 'ClientReceiveMessage'
    | 'ClientCreateChat'
    | 'ClientReceiveChat'
    | 'ClientRequestPrevMessages'
    | 'ClientReceivePrevMessages'

// Base message interface
export interface WebSocketMessage {
    type: MessageType
    data: unknown
}

// Specific message interfaces
export interface ClientSendMessageData {
    chatId: string
    content: string
}

export interface ClientReceiveMessageData {
    chatId: string
    senderName: string
    content: string
    timestamp: string
    messageId: string
}

export interface ClientCreateChatData {
    name: string
    participants: string[]
}

export interface ClientReceiveChatData {
    chatId: string
    name: string
    participants: string[]
}

export interface ClientRequestPrevMessagesData {
    chatId: string
    offset: number
    limit: number
}

export interface ClientReceivePrevMessagesData {
    chatId: string
    messages: Array<{
        messageId: string
        senderName: string
        content: string
        timestamp: string
    }>
    hasMore: boolean
}

// Typed message interfaces
export interface ClientSendMessage extends WebSocketMessage {
    type: 'ClientSendMessage'
    data: ClientSendMessageData
}

export interface ClientReceiveMessage extends WebSocketMessage {
    type: 'ClientReceiveMessage'
    data: ClientReceiveMessageData
}

export interface ClientCreateChat extends WebSocketMessage {
    type: 'ClientCreateChat'
    data: ClientCreateChatData
}

export interface ClientReceiveChat extends WebSocketMessage {
    type: 'ClientReceiveChat'
    data: ClientReceiveChatData
}

export interface ClientRequestPrevMessages extends WebSocketMessage {
    type: 'ClientRequestPrevMessages'
    data: ClientRequestPrevMessagesData
}

export interface ClientReceivePrevMessages extends WebSocketMessage {
    type: 'ClientReceivePrevMessages'
    data: ClientReceivePrevMessagesData
}

// Union type for all message types
export type TypedWebSocketMessage =
    | ClientSendMessage
    | ClientReceiveMessage
    | ClientCreateChat
    | ClientReceiveChat
    | ClientRequestPrevMessages
    | ClientReceivePrevMessages 