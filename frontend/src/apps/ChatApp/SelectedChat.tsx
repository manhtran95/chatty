import MessageList from './MessageList'
import MessageInput from './MessageInput'

interface MessageData {
    senderName: string
    content: string
}

interface SelectedChatProps {
    chatId: string | null
    messages: MessageData[]
}

function SelectedChat({ chatId, messages }: SelectedChatProps) {
    if (chatId === null) {
        return (
            <div className="w-3/4 flex items-center justify-center">
                <div className="text-gray-600">
                    Select a chat to view messages
                </div>
            </div>
        )
    }

    return (
        <div className="w-3/4 h-[80vh] flex flex-col">
            <div className="flex-1 overflow-auto">
                <MessageList messages={messages} />
            </div>
            <MessageInput
                chatId={chatId}
                onMessageSent={() => console.log('Message sent successfully')}
            />
        </div>
    )
}

export default SelectedChat
