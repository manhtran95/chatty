import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import ChatInfoList from './ChatInfoList'
import SelectedChat from './SelectedChat'
import type { ChatData } from './types'
import { receiveNewMessage } from './utils/commandHandlers'
import { useAuth } from '../../modules/auth/AuthContext'
import { useWebSocket } from '../../modules/websocket/WebSocketContext'
import type { WebSocketMessage } from '../../services/WebSocketService'
import './ChatApp.css'
import { stubInitData } from './utils/data'

function ChatApp() {
    const auth = useAuth()
    const navigate = useNavigate()
    const { isConnected, addMessageHandler } = useWebSocket()

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

            // Handle different message types
            switch (message.type) {
                case 'new_message':
                    receiveNewMessage(
                        message.data as {
                            chatId: string
                            senderName: string
                            content: string
                        },
                        setChatListData
                    )
                    break
                case 'chat_update':
                    // Handle chat updates
                    console.log('Chat updated:', message.data)
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

    // User not logged in, return
    if (!isAuthenticated) {
        return (
            <div className="center-container">
                <p>Please log in to access the chat application.</p>
                <button
                    className="btn-primary"
                    onClick={() => navigate('/login')}
                >
                    Go to Login
                </button>
            </div>
        )
    }

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
