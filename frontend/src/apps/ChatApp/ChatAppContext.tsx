import {
    createContext,
    useContext,
    useState,
    useEffect,
    type ReactNode,
} from 'react'
import { commandHandlerReceiveChatList, commandHandlerReceiveNewChat, commandHandlerReceiveNewMessage } from './utils/commandHandlers'
import { useWebSocket } from '../../modules/websocket/WebSocketContext'
import type {
    WebSocketMessage,
    ClientReceiveMessage,
    ClientReceiveChat,
    ClientReceiveChatHistory,
    ClientReceiveChatList
} from '../../services/WebSocketTypes'
import { MESSAGE_TYPES, type ChatInfo, type ChatData } from '../../services/WebSocketTypes'

interface ChatAppContextType {
    // State
    selectedChatId: string | null
    chatListData: ChatData[] | null
    isConnected: boolean

    // Actions
    setSelectedChatId: (chatId: string | null) => void

    // Computed values
    chatInfoList: ChatInfo[]
    selectedChatMessages: Array<{
        senderName: string
        content: string
    }>

    // WebSocket functions
    wsRequestChatList: (offset: number, limit: number) => boolean
    wsRequestChatHistory: (chatId: string, offset: number, limit: number) => boolean
    wsSendMessage: (chatId: string, content: string) => boolean
    wsCreateChat: (name: string, participants: string[]) => boolean
}

const ChatAppContext = createContext<ChatAppContextType | null>(null)

interface ChatAppProviderProps {
    children: ReactNode
}

export const ChatAppProvider = ({ children }: ChatAppProviderProps) => {
    const {
        isConnected,
        addMessageHandler,
        clientRequestChatList: wsRequestChatList,
        clientRequestChatHistory: wsRequestChatHistory,
        clientSendMessage: wsSendMessage,
        clientCreateChat: wsCreateChat
    } = useWebSocket()

    // State
    const [selectedChatId, setSelectedChatId] = useState<string | null>(null)
    const [chatListData, setChatListData] = useState<ChatData[] | null>(null)

    // Request chat list when connected
    useEffect(() => {
        if (!isConnected) return
        wsRequestChatList(0, 20)
    }, [isConnected, wsRequestChatList])

    // Handle WebSocket messages
    useEffect(() => {
        if (!isConnected) return

        const removeHandler = addMessageHandler((message: WebSocketMessage) => {
            console.log('Received WebSocket message:', message)

            // Handle different message types with proper typing
            switch (message.type) {
                case MESSAGE_TYPES.CLIENT_RECEIVE_CHAT_LIST:
                    const receiveChatList = message as ClientReceiveChatList
                    console.log('Received chat list:', receiveChatList.data)
                    commandHandlerReceiveChatList(receiveChatList.data, setChatListData)
                    break
                case MESSAGE_TYPES.CLIENT_RECEIVE_MESSAGE:
                    const receiveMsg = message as ClientReceiveMessage
                    commandHandlerReceiveNewMessage(
                        receiveMsg.data,
                        setChatListData
                    )
                    break
                case MESSAGE_TYPES.CLIENT_RECEIVE_CHAT:
                    const receiveChat = message as ClientReceiveChat
                    console.log('Received new chat:', receiveChat.data)
                    commandHandlerReceiveNewChat(receiveChat.data, setChatListData)
                    break
                case MESSAGE_TYPES.CLIENT_RECEIVE_CHAT_HISTORY:
                    const receiveChatHistory = message as ClientReceiveChatHistory
                    console.log('Received chat history for chat:', receiveChatHistory.data.chatId)
                    console.log('Has more messages:', receiveChatHistory.data.hasMore)
                    // TODO: Handle received chat history
                    break
                default:
                    console.log('Unknown message type:', message.type)
            }
        })

        // Cleanup function to remove the message handler
        return () => {
            removeHandler()
        }
    }, [isConnected, addMessageHandler])

    // Request chat history when chat is selected
    useEffect(() => {
        if (selectedChatId && isConnected) {
            // TODO: Request first 20 messages (offset 0, limit 20)
            // wsRequestChatHistory(selectedChatId, 0, 20)
        }
    }, [selectedChatId, isConnected, wsRequestChatHistory])

    // Computed values
    const chatInfoList: ChatInfo[] = []
    if (chatListData) {
        chatListData.forEach((chat) => {
            chatInfoList.push(chat.chatInfo)
        })
    }

    const selectedChatMessages: Array<{
        senderName: string
        content: string
    }> =
        selectedChatId !== null && chatListData
            ? chatListData.find((chat) => chat.chatInfo.chatId === selectedChatId)
                ?.messages || []
            : []



    const value: ChatAppContextType = {
        // State
        selectedChatId,
        chatListData,
        isConnected,

        // Actions
        setSelectedChatId,

        // Computed values
        chatInfoList,
        selectedChatMessages,

        // WebSocket functions
        wsRequestChatList,
        wsRequestChatHistory,
        wsSendMessage,
        wsCreateChat,
    }

    return (
        <ChatAppContext.Provider value={value}>
            {children}
        </ChatAppContext.Provider>
    )
}

export const useChatApp = () => {
    const context = useContext(ChatAppContext)
    if (!context) {
        throw new Error('useChatApp must be used within a ChatAppProvider')
    }
    return context
} 