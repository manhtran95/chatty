import { chatClasses } from '../../utils/tailwindClasses'

interface ChatInfoProps {
    name: string
    isSelected?: boolean
    onClick?: () => void
}

function ChatInfo({ name, isSelected, onClick }: ChatInfoProps) {
    return (
        <div
            className={isSelected ? chatClasses.itemSelected : chatClasses.itemHover}
            onClick={onClick}
        >
            <div className="font-bold mb-1 text-gray-900">
                {name}
            </div>
        </div>
    )
}

export default ChatInfo
