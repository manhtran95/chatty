import { useState } from 'react'
import { useWebSocket } from '../../modules/websocket/WebSocketContext'
import { buttonClasses, inputClasses } from '../../utils/tailwindClasses'

interface MessageInputProps {
    chatId: string
    onMessageSent?: () => void
}

export default function MessageInput({ chatId, onMessageSent }: MessageInputProps) {
    const [message, setMessage] = useState('')
    const { sendClientMessage, isConnected } = useWebSocket()

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault()

        if (!message.trim() || !isConnected) {
            return
        }

        // Send message via WebSocket using typed function
        const success = sendClientMessage(chatId, message.trim())

        if (success) {
            setMessage('')
            onMessageSent?.()
        } else {
            console.error('Failed to send message')
        }
    }

    return (
        <form onSubmit={handleSubmit} className="flex gap-3 p-4 border-t border-gray-200 bg-gray-50">
            <input
                type="text"
                value={message}
                onChange={(e) => setMessage(e.target.value)}
                placeholder={isConnected ? "Type your message..." : "Connecting..."}
                disabled={!isConnected}
                className={`flex-1 ${isConnected ? inputClasses.base : inputClasses.disabled}`}
            />
            <button
                type="submit"
                disabled={!message.trim() || !isConnected}
                className={buttonClasses.small}
            >
                Send
            </button>
        </form>
    )
}
