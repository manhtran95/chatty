import type { ChatData, NewMessage } from '../types'

export function receiveNewMessage(
    newMessage: NewMessage,
    setChatListData: React.Dispatch<React.SetStateAction<ChatData[] | null>>
) {
    // Update the chat list data with the new message, add it to head of the messages array of the chat
    setChatListData((prevData) => {
        if (!prevData) return null
        return prevData.map((chat) => {
            if (chat.id === newMessage.chatId) {
                return {
                    ...chat,
                    lastMessage: newMessage.content, // update last message
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
