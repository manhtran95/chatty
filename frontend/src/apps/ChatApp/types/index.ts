export interface ChatData {
    chatID: string
    name: string
    participantInfos: Array<{
        id: string
        name: string
        email: string
    }>
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

// export interface NewChat {
//     id: string
//     name: string
//     participantEmails: Array<string>
// }