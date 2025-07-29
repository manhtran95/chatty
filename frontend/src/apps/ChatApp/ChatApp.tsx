import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import ChatInfoList from './ChatInfoList'
import SelectedChat from './SelectedChat'
import type { ChatData } from './types'
import { receiveNewChat, receiveNewMessage } from './utils/commandHandlers'
import { useAuth } from '../../modules/auth/AuthContext'
import { useWebSocket } from '../../modules/websocket/WebSocketContext'
import type {
    WebSocketMessage,
    ClientReceiveMessage,
    ClientReceiveChat,
    ClientReceiveChatHistory
} from '../../services/WebSocketTypes'
import { MESSAGE_TYPES } from '../../services/WebSocketTypes'
import { buttonClasses, statusClasses, layoutClasses } from '../../utils/tailwindClasses'
import './styles/ChatApp.css'
import { stubInitData } from './utils/data'

function ChatApp() {
    const auth = useAuth()
    const navigate = useNavigate()
    const { isConnected, addMessageHandler, requestChatHistory } = useWebSocket()

    const { isAuthenticated, user, logout } = auth!
    console.log(`user: ${user}`)
    // states
    const [selectedChatId, setSelectedChatId] = useState<string | null>(null)
    const [chatListData, setChatListData] = useState<ChatData[] | null>(null)
    const [showUserDropdown, setShowUserDropdown] = useState(false)

    // run once when the component mounts
    useEffect(() => {
        setChatListData(stubInitData)
    }, []) // Empty dependency array = run once

    // Handle WebSocket messages
    useEffect(() => {
        if (!isConnected) return

        const removeHandler = addMessageHandler((message: WebSocketMessage) => {
            console.log('Received WebSocket message:', message)

            // Handle different message types with proper typing
            switch (message.type) {
                case MESSAGE_TYPES.CLIENT_RECEIVE_MESSAGE:
                    const receiveMsg = message as ClientReceiveMessage
                    receiveNewMessage(
                        receiveMsg.data,
                        setChatListData
                    )
                    break
                case MESSAGE_TYPES.CLIENT_RECEIVE_CHAT:
                    const receiveChat = message as ClientReceiveChat
                    console.log('Received new chat:', receiveChat.data)
                    receiveNewChat(receiveChat.data, setChatListData)
                    break
                case MESSAGE_TYPES.CLIENT_RECEIVE_CHAT_HISTORY:
                    const receiveChatHistory = message as ClientReceiveChatHistory
                    console.log('Received chat history for chat:', receiveChatHistory.data.chatID)
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
            // requestChatHistory(selectedChatId, 0, 20)
        }
    }, [selectedChatId, isConnected, requestChatHistory])

    // User not logged in, return
    // if (!isAuthenticated) {
    //     return (
    //         <div className="center-container">
    //             <p>Please log in to access the chat application.</p>
    //             <button
    //                 className={buttonClasses.primary}
    //                 onClick={() => navigate('/login')}
    //             >
    //                 Go to Login
    //             </button>
    //         </div>
    //     )
    // }

    // deduced data
    const chatInfoList = chatListData
        ? chatListData.map((chat) => ({
            chatID: chat.chatID,
            name: chat.name,
        }))
        : []
    const selectedChatMessages: Array<{
        senderName: string
        content: string
    }> =
        selectedChatId !== null && chatListData
            ? chatListData.find((chat) => chat.chatID === selectedChatId)
                ?.messages || []
            : []

    const handleLogout = () => {
        logout()
        navigate('/login')
    }

    return (
        <>
            {/* Top Bar */}
            <div className={`${layoutClasses.flexBetween} p-4 bg-gray-50 border-b border-gray-200`}>
                <div className="relative">
                    <button
                        onClick={() => setShowUserDropdown(!showUserDropdown)}
                        className="flex items-center gap-2 px-3 py-2 bg-white border border-gray-300 rounded-md text-sm text-gray-700 hover:bg-gray-50 transition-colors duration-150"
                        onBlur={() =>
                            setTimeout(() => setShowUserDropdown(false), 150)
                        }
                    >
                        <span>Hi, {user?.name || 'User'}</span>
                        <span className="text-xs">▼</span>
                    </button>

                    {showUserDropdown && (
                        <div className="absolute top-full right-0 mt-1 bg-white border border-gray-200 rounded-md shadow-lg z-50 min-w-32">
                            <button
                                onClick={handleLogout}
                                className="w-full px-4 py-2 text-left text-sm text-red-600 hover:bg-gray-50"
                            >
                                Logout
                            </button>
                        </div>
                    )}
                </div>
                {/* WebSocket connection status */}
                <div className="flex items-center">
                    <span className={isConnected ? statusClasses.connected : statusClasses.disconnected}>
                        {isConnected ? '🟢 Connected' : '🔴 Disconnected'}
                    </span>
                </div>
            </div>

            <div className="flex h-[80vh]">
                <ChatInfoList
                    chats={chatInfoList}
                    selectedChatId={selectedChatId}
                    onChatSelect={setSelectedChatId}
                />
                <SelectedChat
                    chatId={selectedChatId}
                    messages={selectedChatMessages}
                />
            </div>
            {/* <div className="p-4">
                <button
                    onClick={() =>
                        receiveNewMessage(
                            {
                                chatID: '1',
                                senderName: 'manh',
                                content: `New message here!${Math.random()}`,
                                timestamp: new Date().toISOString(),
                                messageId: '1',
                            },
                            setChatListData
                        )
                    }
                    className={buttonClasses.small}
                >
                    Receive new mess testing
                </button>
            </div> */}
        </>
    )
}

export default ChatApp
