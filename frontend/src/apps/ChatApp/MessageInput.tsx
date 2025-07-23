import { useState } from 'react'
import { useWebSocket } from '../../modules/websocket/WebSocketContext'

interface MessageInputProps {
    chatId: string
    onMessageSent?: () => void
}

export default function MessageInput({
    chatId,
    onMessageSent,
}: MessageInputProps) {
    const [message, setMessage] = useState('')
    const { sendMessage, isConnected } = useWebSocket()

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault()

        if (!message.trim() || !isConnected) {
            return
        }

        // Send message via WebSocket
        const success = sendMessage({
            type: 'send_message',
            data: {
                chatId,
                content: message.trim(),
                timestamp: new Date().toISOString(),
            },
        })

        if (success) {
            setMessage('')
            onMessageSent?.()
        } else {
            console.error('Failed to send message')
        }
    }

    return (
        <form onSubmit={handleSubmit} className="message-input-form">
            <input
                type="text"
                value={message}
                onChange={(e) => setMessage(e.target.value)}
                placeholder={
                    isConnected ? 'Type your message...' : 'Connecting...'
                }
                disabled={!isConnected}
                className="message-input"
            />
            <button
                type="submit"
                disabled={!message.trim() || !isConnected}
                className="send-button"
            >
                Send
            </button>
        </form>
    )
}
