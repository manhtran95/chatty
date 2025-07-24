import { chatClasses } from '../../utils/tailwindClasses'

interface ChatInfoProps {
    name: string
    lastMessage: string
    isSelected?: boolean
    onClick?: () => void
}

function ChatInfo({ name, lastMessage, isSelected, onClick }: ChatInfoProps) {
    return (
        <div
            className={isSelected ? chatClasses.itemSelected : chatClasses.itemHover}
            onClick={onClick}
        >
            <div className="font-bold mb-1 text-gray-900">
                {name}
            </div>
            <div className="text-gray-600 text-sm truncate">
                {lastMessage}
            </div>
        </div>
    )
}

export default ChatInfo
