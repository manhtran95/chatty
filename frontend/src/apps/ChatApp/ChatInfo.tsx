interface ChatInfoProps {
    name: string
    lastMessage: string
    isSelected?: boolean
    onClick?: () => void
}

function ChatInfo({ name, lastMessage, isSelected, onClick }: ChatInfoProps) {
    return (
        <div
            style={{
                padding: '12px',
                borderBottom: '1px solid #eee',
                cursor: 'pointer',
                backgroundColor: isSelected ? '#f0f8ff' : 'white',
                textAlign: 'left',
            }}
            onClick={onClick}
        >
            <div style={{ fontWeight: 'bold', marginBottom: '4px' }}>
                {name}
            </div>
            <div style={{ color: '#666', fontSize: '14px' }}>{lastMessage}</div>
        </div>
    )
}

export default ChatInfo
