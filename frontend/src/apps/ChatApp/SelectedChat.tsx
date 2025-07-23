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
            <div
                style={{
                    width: '75%',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                }}
            >
                <div style={{ color: '#666' }}>
                    Select a chat to view messages
                </div>
            </div>
        )
    }

    return (
        <div
            style={{
                width: '75%',
                height: '80vh',
                display: 'flex',
                flexDirection: 'column',
            }}
        >
            <div style={{ flex: 1, overflow: 'auto' }}>
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
