// WebSocketTypes.ts

// Message type constants
export const MESSAGE_TYPES = {
    CLIENT_SEND_MESSAGE: 'ClientSendMessage',
    CLIENT_RECEIVE_MESSAGE: 'ClientReceiveMessage',
    CLIENT_CREATE_CHAT: 'ClientCreateChat',
    CLIENT_RECEIVE_CHAT: 'ClientReceiveChat',
    CLIENT_REQUEST_CHAT_HISTORY: 'ClientRequestChatHistory',
    CLIENT_RECEIVE_CHAT_HISTORY: 'ClientReceiveChatHistory',
} as const

// Message type definitions
export type MessageType = typeof MESSAGE_TYPES[keyof typeof MESSAGE_TYPES]

// Base message interface
export interface WebSocketMessage {
    type: MessageType
    data: unknown
    senderId: string
}

export interface WebSocketMessageResponse {
    type: MessageType
    success: boolean
    error: string
}

// Specific message interfaces
export interface ClientSendMessageData {
    chatID: string
    senderId: string
    content: string
}

export interface ClientReceiveMessageData {
    chatID: string
    senderName: string
    content: string
    timestamp: string
    messageId: string
}

export interface ClientCreateChatData {
    name: string
    participantEmails: string[]
}

export interface ClientReceiveChatData {
    chatID: string
    name: string
    participantInfos: Array<{
        id: string
        name: string
        email: string
    }>
}

export interface ClientRequestChatHistoryData {
    chatID: string
    offset: number
    limit: number
}

export interface ClientReceiveChatHistoryData {
    chatID: string
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
    type: typeof MESSAGE_TYPES.CLIENT_SEND_MESSAGE
    data: ClientSendMessageData
}

export interface ClientReceiveMessage extends WebSocketMessage {
    type: typeof MESSAGE_TYPES.CLIENT_RECEIVE_MESSAGE
    data: ClientReceiveMessageData
}

export interface ClientCreateChat extends WebSocketMessage {
    type: typeof MESSAGE_TYPES.CLIENT_CREATE_CHAT
    data: ClientCreateChatData
}

export interface ClientReceiveChat extends WebSocketMessage {
    type: typeof MESSAGE_TYPES.CLIENT_RECEIVE_CHAT
    data: ClientReceiveChatData
}

export interface ClientRequestChatHistory extends WebSocketMessage {
    type: typeof MESSAGE_TYPES.CLIENT_REQUEST_CHAT_HISTORY
    data: ClientRequestChatHistoryData
}

export interface ClientReceiveChatHistory extends WebSocketMessage {
    type: typeof MESSAGE_TYPES.CLIENT_RECEIVE_CHAT_HISTORY
    data: ClientReceiveChatHistoryData
}

// Union type for all message types
export type TypedWebSocketMessage =
    | ClientSendMessage
    | ClientReceiveMessage
    | ClientCreateChat
    | ClientReceiveChat
    | ClientRequestChatHistory
    | ClientReceiveChatHistory 