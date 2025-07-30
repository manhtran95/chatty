export interface ChatInfoType {
    chatId: string
    name: string
    participantInfos: Array<{
        id: string
        name: string
        email: string
    }>
}
export interface ChatData {
    chatInfo: ChatInfoType
    // order of messages: latest message first
    messages: Array<{
        senderName: string
        content: string
    }>
}

export interface NewMessage {
    chatID: string
    senderName: string
    content: string
}
