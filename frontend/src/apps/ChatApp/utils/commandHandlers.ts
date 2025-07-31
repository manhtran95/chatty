import type { ChatData, ClientReceiveChatData, ClientReceiveChatListData, ClientReceiveMessageData } from '../../../services/WebSocketTypes'
// import type { ChatData } from '../types'

// 1
export function commandHandlerReceiveChatList(
    chatList: ClientReceiveChatListData,
    setChatListData: React.Dispatch<React.SetStateAction<ChatData[] | null>>
) {
    setChatListData(chatList.chats.map((chat) => ({
        chatInfo: {
            chatId: chat.chatId,
            name: chat.name,
            participantInfos: chat.participantInfos,
        },
        messages: [],
    })))
}

// 2
export function commandHandlerReceiveNewChat(
    newChat: ClientReceiveChatData,
    setChatListData: React.Dispatch<React.SetStateAction<ChatData[] | null>>
) {
    // Update the chat list data with the new chat, add it to head of the chats array
    setChatListData((prevData) => {
        if (!prevData) return null
        return [
            {
                chatInfo: {
                    chatId: newChat.chatId,
                    name: newChat.name,
                    participantInfos: newChat.participantInfos,
                },
                messages: [],
            },
            ...prevData,
        ]
    })
}

// 4
export function commandHandlerReceiveNewMessage(
    newMessage: ClientReceiveMessageData,
    setChatListData: React.Dispatch<React.SetStateAction<ChatData[] | null>>
) {
    // Update the chat list data with the new message, add it to head of the messages array of the chat
    setChatListData((prevData) => {
        if (!prevData) return null
        return prevData.map((chat) => {
            if (chat.chatInfo.chatId === newMessage.chatId) {
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