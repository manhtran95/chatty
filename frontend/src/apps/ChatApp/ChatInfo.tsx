import { chatClasses } from '../../utils/tailwindClasses'

interface ChatInfoProps {
    name: string
    participantInfos: Array<{
        id: string
        name: string
        email: string
    }>
    isSelected?: boolean
    onClick?: () => void
}

function ChatInfo({ name, participantInfos, isSelected, onClick }: ChatInfoProps) {
    return (
        <div
            className={isSelected ? chatClasses.itemSelected : chatClasses.itemHover}
            onClick={onClick}
        >
            <div className="font-bold mb-1 text-gray-900">
                {name}
            </div>
            <div className="text-gray-500 text-sm">
                {participantInfos.map((participant) => participant.email).join(', ')}
            </div>
        </div>
    )
}

export default ChatInfo
