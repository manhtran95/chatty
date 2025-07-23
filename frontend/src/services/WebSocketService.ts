// WebSocketService.ts
let socket: WebSocket | null = null
let isConnected = false
let messageHandlers: Array<(data: WebSocketMessage) => void> = []

export interface WebSocketMessage {
    type: string
    data: unknown
}

export function connectWebSocket(token: string): Promise<boolean> {
    return new Promise((resolve) => {
        if (socket && socket.readyState === WebSocket.OPEN) {
            console.log('WebSocket already connected')
            resolve(true)
            return
        }

        socket = new WebSocket(`ws://localhost:8080/ws?token=${token}`)

        socket.onopen = () => {
            console.log('WebSocket connected')
            isConnected = true
            resolve(true)
        }

        socket.onmessage = (event) => {
            try {
                const message: WebSocketMessage = JSON.parse(event.data)
                console.log('Message from server:', message)

                // Notify all registered handlers
                messageHandlers.forEach((handler) => {
                    try {
                        handler(message)
                    } catch (error) {
                        console.error('Error in message handler:', error)
                    }
                })
            } catch (error) {
                console.error('Error parsing WebSocket message:', error)
            }
        }

        socket.onclose = (event) => {
            console.log('WebSocket disconnected', event.code, event.reason)
            isConnected = false
            socket = null
        }

        socket.onerror = (error) => {
            console.error('WebSocket error:', error)
            isConnected = false
            resolve(false)
        }
    })
}

export function sendMessage(msg: WebSocketMessage): boolean {
    if (socket && socket.readyState === WebSocket.OPEN) {
        socket.send(JSON.stringify(msg))
        return true
    } else {
        console.warn('WebSocket not connected. Message not sent.')
        return false
    }
}

export function disconnectWebSocket(): boolean {
    if (socket && socket.readyState === WebSocket.OPEN) {
        socket.close(1000, 'Client closed connection')
        isConnected = false
        socket = null
        return true
    }
    return false
}

export function isWebSocketConnected(): boolean {
    return isConnected && socket?.readyState === WebSocket.OPEN
}

export function addMessageHandler(
    handler: (data: WebSocketMessage) => void
): () => void {
    messageHandlers.push(handler)

    // Return a function to remove this handler
    return () => {
        const index = messageHandlers.indexOf(handler)
        if (index > -1) {
            messageHandlers.splice(index, 1)
        }
    }
}

export function clearMessageHandlers(): void {
    messageHandlers = []
}
