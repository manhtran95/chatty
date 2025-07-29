import type { ClientReceiveChatData, ClientReceiveMessageData } from '../../../services/WebSocketTypes'
import type { ChatData } from '../types'

export function receiveNewMessage(
    newMessage: ClientReceiveMessageData,
    setChatListData: React.Dispatch<React.SetStateAction<ChatData[] | null>>
) {
    // Update the chat list data with the new message, add it to head of the messages array of the chat
    setChatListData((prevData) => {
        if (!prevData) return null
        return prevData.map((chat) => {
            if (chat.chatID === newMessage.chatID) {
                return {
                    ...chat,
                    messages: [
                        {
                            senderName: newMessage.senderName,
                            content: newMessage.content,
                        },
                        ...chat.messages, // prepend the new message to the existing messages
                    ],
                }
            }
            return chat
        })
    })
}

export function receiveNewChat(
    newChat: ClientReceiveChatData,
    setChatListData: React.Dispatch<React.SetStateAction<ChatData[] | null>>
) {
    // Update the chat list data with the new chat, add it to head of the chats array
    setChatListData((prevData) => {
        if (!prevData) return null
        return [
            {
                chatID: newChat.chatID,
                name: newChat.name,
                participantInfos: newChat.participantInfos,
                messages: [],
            },
            ...prevData,
        ]
    })
}
