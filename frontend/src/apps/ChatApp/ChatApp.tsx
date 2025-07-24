import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import ChatInfoList from './ChatInfoList'
import SelectedChat from './SelectedChat'
import type { ChatData } from './types'
import { receiveNewMessage } from './utils/commandHandlers'
import { useAuth } from '../../modules/auth/AuthContext'
import { useWebSocket } from '../../modules/websocket/WebSocketContext'
import type {
    WebSocketMessage,
    ClientReceiveMessage,
    ClientReceiveChat,
    ClientReceivePrevMessages
} from '../../services/WebSocketTypes'
import './styles/ChatApp.css'
import { stubInitData } from './utils/data'

function ChatApp() {
    const auth = useAuth()
    const navigate = useNavigate()
    const { isConnected, addMessageHandler, requestPrevMessages } = useWebSocket()

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
                case 'ClientReceiveMessage':
                    const receiveMsg = message as ClientReceiveMessage
                    receiveNewMessage(
                        {
                            chatId: receiveMsg.data.chatId,
                            senderName: receiveMsg.data.senderName,
                            content: receiveMsg.data.content,
                        },
                        setChatListData
                    )
                    break
                case 'ClientReceiveChat':
                    const receiveChat = message as ClientReceiveChat
                    console.log('Received new chat:', receiveChat.data)
                    // TODO: Handle new chat creation
                    break
                case 'ClientReceivePrevMessages':
                    const receivePrevMsg = message as ClientReceivePrevMessages
                    console.log('Received previous messages for chat:', receivePrevMsg.data.chatId)
                    console.log('Has more messages:', receivePrevMsg.data.hasMore)
                    // TODO: Handle received previous messages
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

    // Request previous messages when chat is selected
    useEffect(() => {
        if (selectedChatId && isConnected) {
            // Request first 20 messages (offset 0, limit 20)
            requestPrevMessages(selectedChatId, 0, 20)
        }
    }, [selectedChatId, isConnected, requestPrevMessages])

    // User not logged in, return
    // if (!isAuthenticated) {
    //     return (
    //         <div className="center-container">
    //             <p>Please log in to access the chat application.</p>
    //             <button
    //                 className="btn-primary"
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
            id: chat.id,
            name: chat.name,
            lastMessage: chat.lastMessage || 'No messages',
        }))
        : []
    const selectedChatMessages: Array<{
        senderName: string
        content: string
    }> =
        selectedChatId !== null && chatListData
            ? chatListData.find((chat) => chat.id === selectedChatId)
                ?.messages || []
            : []

    const handleLogout = () => {
        logout()
        navigate('/login')
    }

    return (
        <>
            {/* Top Bar */}
            <div className="top-bar">
                <div className="user-dropdown-container">
                    <button
                        onClick={() => setShowUserDropdown(!showUserDropdown)}
                        className="user-dropdown-button"
                        onBlur={() =>
                            setTimeout(() => setShowUserDropdown(false), 150)
                        }
                    >
                        <span>Hi, {user?.name || 'User'}</span>
                        <span className="dropdown-arrow">▼</span>
                    </button>

                    {showUserDropdown && (
                        <div className="user-dropdown-menu">
                            <button
                                onClick={handleLogout}
                                className="user-dropdown-item"
                            >
                                Logout
                            </button>
                        </div>
                    )}
                </div>
                {/* WebSocket connection status */}
                <div className="connection-status">
                    <span
                        className={`status-indicator ${isConnected ? 'connected' : 'disconnected'}`}
                    >
                        {isConnected ? '🟢 Connected' : '🔴 Disconnected'}
                    </span>
                </div>
            </div>

            <div className="chat-container">
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
            <div>
                <button
                    onClick={() =>
                        receiveNewMessage(
                            {
                                chatId: '1',
                                senderName: 'manh',
                                content: `New message here!${Math.random()}`,
                            },
                            setChatListData
                        )
                    }
                >
                    Receive new mess testing
                </button>
            </div>
        </>
    )
}

export default ChatApp
