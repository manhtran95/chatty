export interface ChatData {
    id: string;
    name: string;
    userList: Array<{
        id: string;
        username: string;
    }>;
    lastMessage: string;
    // order of messages: latest message first
    messages: Array<{
        senderName: string;
        content: string;
    }>;
}

export interface NewMessage {
    chatId: string;
    senderName: string;
    content: string;
}