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
    type WebSocketMessage,
} from '../../services/WebSocketService'
import { useAuth } from '../auth/AuthContext'

interface WebSocketContextType {
    isConnected: boolean
    connect: () => Promise<boolean>
    disconnect: () => boolean
    sendMessage: (message: WebSocketMessage) => boolean
    addMessageHandler: (handler: (data: WebSocketMessage) => void) => () => void
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
