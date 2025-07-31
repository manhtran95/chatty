// WebSocketTypes.ts
// Message type constants
export const MESSAGE_TYPES = {
    // request command
    CLIENT_REQUEST_CHAT_LIST: 'ClientRequestChatList', // A.1
    CLIENT_CREATE_CHAT: 'ClientCreateChat', // A.2
    CLIENT_REQUEST_CHAT_HISTORY: 'ClientRequestChatHistory', // A.3
    CLIENT_SEND_MESSAGE: 'ClientSendMessage', // A.4
    // response command
    CLIENT_RECEIVE_CHAT_LIST: 'ClientReceiveChatList', // B.1
    CLIENT_RECEIVE_CHAT: 'ClientReceiveChat', // B.2
    CLIENT_RECEIVE_CHAT_HISTORY: 'ClientReceiveChatHistory', // B.3
    CLIENT_RECEIVE_MESSAGE: 'ClientReceiveMessage', // B.4
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

export interface ChatInfo {
    chatId: string
    name: string
    participantInfos: Array<{
        id: string
        name: string
        email: string
    }>
}

export interface ChatData {
    chatInfo: ChatInfo
    // order of messages: latest message first
    messages: Array<{
        senderName: string
        content: string
    }>
}

// ***
// 1. CHAT LIST
// ***
export interface ClientRequestChatListData {
    offset: number
    limit: number
}

export interface ClientReceiveChatListData {
    chats: Array<ChatInfo>
}

// ***
// 2. CREATE CHAT
// ***
export interface ClientCreateChatData {
    name: string
    participantEmails: string[]
}

export type ClientReceiveChatData = ChatInfo

// ***
// 3. CHAT HISTORY
// ***
export interface ClientRequestChatHistoryData {
    chatId: string
    offset: number
    limit: number
}

export interface ClientReceiveChatHistoryData {
    chatId: string
    messages: Array<{
        messageId: string
        senderName: string
        content: string
        timestamp: string
    }>
    hasMore: boolean
}


// ***
// 4. NEW MESSAGE
// ***
export interface ClientSendMessageData {
    chatId: string
    senderId: string
    content: string
}

export interface ClientReceiveMessageData {
    chatId: string
    senderName: string
    content: string
    timestamp: string
    messageId: string
}

// ***
// ***
// Specific message interfaces
// ***
// ***

// 1. CHAT LIST
export interface ClientRequestChatList extends WebSocketMessage {
    type: typeof MESSAGE_TYPES.CLIENT_REQUEST_CHAT_LIST
    data: ClientRequestChatListData
}

export interface ClientReceiveChatList extends WebSocketMessage {
    type: typeof MESSAGE_TYPES.CLIENT_RECEIVE_CHAT_LIST
    data: ClientReceiveChatListData
}

// 2. NEW MESSAGE
export interface ClientSendMessage extends WebSocketMessage {
    type: typeof MESSAGE_TYPES.CLIENT_SEND_MESSAGE
    data: ClientSendMessageData
}

export interface ClientReceiveMessage extends WebSocketMessage {
    type: typeof MESSAGE_TYPES.CLIENT_RECEIVE_MESSAGE
    data: ClientReceiveMessageData
}

// 2. CREATE CHAT
export interface ClientCreateChat extends WebSocketMessage {
    type: typeof MESSAGE_TYPES.CLIENT_CREATE_CHAT
    data: ClientCreateChatData
}

export interface ClientReceiveChat extends WebSocketMessage {
    type: typeof MESSAGE_TYPES.CLIENT_RECEIVE_CHAT
    data: ClientReceiveChatData
}

// 3. CHAT HISTORY
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