import { useState } from 'react'
import ChatInfo from './ChatInfo'
import CreateChatModal from './CreateChatModal'
import { buttonClasses, chatClasses } from '../../utils/tailwindClasses'
import type { ChatInfoType } from './types'

// interface ChatData {
//     chatID: string
//     name: string
// }

interface ChatInfoListProps {
    chatInfoList: Array<ChatInfoType>
    selectedChatId: string | null
    onChatSelect: (chatId: string) => void
}

function ChatInfoList({
    chatInfoList,
    selectedChatId,
    onChatSelect,
}: ChatInfoListProps) {
    const [isModalOpen, setIsModalOpen] = useState(false)

    const handleCreateChat = () => {
        setIsModalOpen(true)
    }

    const handleCloseModal = () => {
        setIsModalOpen(false)
    }

    console.log('ChatInfoList chatInfoList', chatInfoList)

    return (
        <div className="w-1/4 border-r border-gray-300 h-[80vh] flex flex-col">
            {/* Create Chat Button */}
            <div className={chatClasses.header}>
                <button
                    onClick={handleCreateChat}
                    className="w-full px-4 py-3 bg-blue-600 text-white border-none rounded-md text-sm font-medium cursor-pointer flex items-center justify-center gap-2 transition-colors duration-150 hover:bg-blue-700 active:bg-blue-800"
                    title="Create new chat"
                >
                    <span className="text-lg font-bold leading-none">+</span>
                    <span>New Chat</span>
                </button>
            </div>

            {/* Chat List */}
            <div className={chatClasses.list}>
                {chatInfoList.map((chat) => (
                    <ChatInfo
                        key={chat.chatId}
                        name={chat.name}
                        participantInfos={chat.participantInfos}
                        isSelected={selectedChatId === chat.chatId}
                        onClick={() => onChatSelect(chat.chatId)}
                    />
                ))}
            </div>

            {/* Create Chat Modal */}
            <CreateChatModal isOpen={isModalOpen} onClose={handleCloseModal} />
        </div>
    )
}

export default ChatInfoList
