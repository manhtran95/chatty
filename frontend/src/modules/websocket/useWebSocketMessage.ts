import { useEffect } from 'react'
import { useWebSocket } from './WebSocketContext'
import type { WebSocketMessage } from '../../services/WebSocketService'

interface UseWebSocketMessageOptions {
    onMessage?: (data: WebSocketMessage) => void
    onConnect?: () => void
    onDisconnect?: () => void
}

export function useWebSocketMessage(options: UseWebSocketMessageOptions = {}) {
    const { isConnected, addMessageHandler } = useWebSocket()

    useEffect(() => {
        if (!isConnected) {
            options.onDisconnect?.()
            return
        }

        options.onConnect?.()

        const removeHandler = addMessageHandler((message) => {
            options.onMessage?.(message)
        })

        return () => {
            removeHandler()
        }
    }, [isConnected, addMessageHandler, options])

    return { isConnected }
}
