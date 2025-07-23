import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import ChatInfoList from './ChatInfoList'
import SelectedChat from './SelectedChat'
import type { ChatData } from './types'
import { receiveNewMessage } from './utils/commandHandlers'
import { useAuth } from '../../modules/auth/AuthContext'
import './ChatApp.css'
import { stubInitData } from './utils/data'

function ChatApp() {
    const auth = useAuth()
    const navigate = useNavigate()

    const { isAuthenticated, user, logout } = auth!
    console.log(`user: ${user}`)
    // states
    const [selectedChatId, setSelectedChatId] = useState<string | null>(null)
    const [chatListData, setChatListData] = useState<ChatData[] | null>(null)
    const [showUserDropdown, setShowUserDropdown] = useState(false)

    // stub data


    console.log(`selectedChatId: ${selectedChatId}`)

    // run once when the component mounts
    useEffect(() => {
        setChatListData(stubInitData)
        // Optional cleanup
        // return () => {
        //     console.log('Component unmounted')
        // }
    }, []) // Empty dependency array = run once

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

    return (
        <>
            {/* Top Bar */}
            <div className="top-bar">
                <div className="user-dropdown-container">
                    <button
                        onClick={() => setShowUserDropdown(!showUserDropdown)}
                        className="user-dropdown-button"
                        onBlur={() => setTimeout(() => setShowUserDropdown(false), 150)}
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
