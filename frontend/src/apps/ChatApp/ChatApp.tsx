import { useChatApp } from './ChatAppContext'
import ChatInfoList from './ChatInfoList'
import SelectedChat from './SelectedChat'
import { useAuth } from '../../modules/auth/AuthContext'
import { statusClasses, layoutClasses } from '../../utils/tailwindClasses'
import './styles/ChatApp.css'
import { useState } from 'react'
import { useNavigate } from 'react-router-dom'

function ChatApp() {
    const navigate = useNavigate()
    
    const auth = useAuth()
    const { user, logout } = auth!
    
    const {
        selectedChatId,
        isConnected,
        chatInfoList,
        selectedChatMessages,
        setSelectedChatId,
    } = useChatApp()

    const [showUserDropdown, setShowUserDropdown] = useState(false)
    
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
                    chatInfoList={chatInfoList}
                    selectedChatId={selectedChatId}
                    onChatSelect={setSelectedChatId}
                />
                <SelectedChat
                    chatId={selectedChatId}
                    messages={selectedChatMessages}
                />
            </div>
        </>
    )
}

export default ChatApp
