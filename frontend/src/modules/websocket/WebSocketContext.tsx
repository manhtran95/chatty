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
    ClientRequestChatHistory,
    ClientReceiveChatHistory,
    ClientRequestChatList,
} from '../../services/WebSocketTypes'
import { MESSAGE_TYPES } from '../../services/WebSocketTypes'
import { useAuth } from '../auth/AuthContext'

interface WebSocketContextType {
    isConnected: boolean
    connect: () => Promise<boolean>
    disconnect: () => boolean
    sendMessage: (message: WebSocketMessage) => boolean
    addMessageHandler: (handler: (data: WebSocketMessage) => void) => () => void
    // command functions
    clientSendMessage: (chatId: string, content: string) => boolean
    clientCreateChat: (name: string, participants: string[]) => boolean
    clientRequestChatHistory: (chatId: string, offset: number, limit: number) => boolean
    clientRequestChatList: (offset: number, limit: number) => boolean
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

    // 1. CHAT LIST
    const requestChatList = (offset: number, limit: number): boolean => {
        const message: ClientRequestChatList = {
            type: MESSAGE_TYPES.CLIENT_REQUEST_CHAT_LIST,
            data: { offset, limit },
            senderId: auth?.user?.id || '',
        }
        return sendMessage(message)
    }

    // 4. NEW MESSAGE
    const sendClientMessage = (chatId: string, content: string): boolean => {
        const message: ClientSendMessage = {
            type: MESSAGE_TYPES.CLIENT_SEND_MESSAGE,
            data: {
                chatId,
                content,
                senderId: auth?.user?.id || '',
            },
            senderId: auth?.user?.id || '',
        }
        return sendMessage(message)
    }

    // 2. CREATE CHAT
    const createChat = (name: string, participantEmails: string[]): boolean => {
        participantEmails.push(auth?.user?.email || '')
        const message: ClientCreateChat = {
            type: MESSAGE_TYPES.CLIENT_CREATE_CHAT,
            data: {
                name,
                participantEmails: participantEmails,
            },
            senderId: auth?.user?.id || '',
        }
        return sendMessage(message)
    }

    // 3. CHAT HISTORY
    const requestChatHistory = (chatId: string, offset: number, limit: number): boolean => {
        const message: ClientRequestChatHistory = {
            type: MESSAGE_TYPES.CLIENT_REQUEST_CHAT_HISTORY,
            data: {
                chatId,
                offset,
                limit,
            },
            senderId: auth?.user?.id || '',
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
            // disconnect()
        }
    }, [auth?.isAuthenticated, auth?.accessToken])

    // Cleanup on unmount
    useEffect(() => {
        console.log('WebSocketProvider unmount')
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
                // command functions
                clientRequestChatList: requestChatList,
                clientCreateChat: createChat,
                clientRequestChatHistory: requestChatHistory,
                clientSendMessage: sendClientMessage,
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
