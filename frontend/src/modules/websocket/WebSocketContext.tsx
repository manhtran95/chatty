import {
    createContext,
    useContext,
    useState,
    useEffect,
    type ReactNode,
} from 'react'
import {
    connectWebSocket,
    disconnectWebSocket,
    sendMessage,
    isWebSocketConnected,
    addMessageHandler,
    clearMessageHandlers,
} from '../../services/WebSocketService'
import type {
    WebSocketMessage,
    TypedWebSocketMessage,
    ClientSendMessage,
    ClientReceiveMessage,
    ClientCreateChat,
    ClientReceiveChat,
    ClientRequestPrevMessages,
    ClientReceivePrevMessages,
} from '../../services/WebSocketTypes'
import { useAuth } from '../auth/AuthContext'

interface WebSocketContextType {
    isConnected: boolean
    connect: () => Promise<boolean>
    disconnect: () => boolean
    sendMessage: (message: WebSocketMessage) => boolean
    addMessageHandler: (handler: (data: WebSocketMessage) => void) => () => void
    // Typed message sending functions
    sendClientMessage: (chatId: string, content: string) => boolean
    createChat: (name: string, participants: string[]) => boolean
    requestPrevMessages: (chatId: string, offset: number, limit: number) => boolean
}

const WebSocketContext = createContext<WebSocketContextType | null>(null)

interface WebSocketProviderProps {
    children: ReactNode
}

export const WebSocketProvider = ({ children }: WebSocketProviderProps) => {
    const [isConnected, setIsConnected] = useState(false)
    const auth = useAuth()

    const connect = async (): Promise<boolean> => {
        if (!auth?.accessToken) {
            console.warn('No access token available for WebSocket connection')
            return false
        }

        try {
            const success = await connectWebSocket(auth.accessToken)
            setIsConnected(success)
            return success
        } catch (error) {
            console.error('Failed to connect WebSocket:', error)
            setIsConnected(false)
            return false
        }
    }

    const disconnect = (): boolean => {
        const success = disconnectWebSocket()
        setIsConnected(false)
        return success
    }

    const sendMessageToServer = (message: WebSocketMessage): boolean => {
        return sendMessage(message)
    }

    const addMessageHandlerToContext = (
        handler: (data: WebSocketMessage) => void
    ) => {
        return addMessageHandler(handler)
    }

    // Typed message sending functions
    const sendClientMessage = (chatId: string, content: string): boolean => {
        const message: ClientSendMessage = {
            type: 'ClientSendMessage',
            data: {
                chatId,
                content,
            },
        }
        return sendMessage(message)
    }

    const createChat = (name: string, participants: string[]): boolean => {
        const message: ClientCreateChat = {
            type: 'ClientCreateChat',
            data: {
                name,
                participants,
            },
        }
        return sendMessage(message)
    }

    const requestPrevMessages = (chatId: string, offset: number, limit: number): boolean => {
        const message: ClientRequestPrevMessages = {
            type: 'ClientRequestPrevMessages',
            data: {
                chatId,
                offset,
                limit,
            },
        }
        return sendMessage(message)
    }

    // Auto-connect when user is authenticated
    useEffect(() => {
        if (
            auth?.isAuthenticated &&
            auth?.accessToken &&
            !isWebSocketConnected()
        ) {
            connect()
        } else if (!auth?.isAuthenticated) {
            disconnect()
        }
    }, [auth?.isAuthenticated, auth?.accessToken])

    // Cleanup on unmount
    useEffect(() => {
        return () => {
            disconnect()
            clearMessageHandlers()
        }
    }, [])

    return (
        <WebSocketContext.Provider
            value={{
                isConnected,
                connect,
                disconnect,
                sendMessage: sendMessageToServer,
                addMessageHandler: addMessageHandlerToContext,
                sendClientMessage,
                createChat,
                requestPrevMessages,
            }}
        >
            {children}
        </WebSocketContext.Provider>
    )
}

export const useWebSocket = () => {
    const context = useContext(WebSocketContext)
    if (!context) {
        throw new Error('useWebSocket must be used within a WebSocketProvider')
    }
    return context
}
