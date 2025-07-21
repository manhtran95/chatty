import './ChatApp.css'

interface MessageProps {
    senderName: string
    content: string
}

function Message({ senderName, content }: MessageProps) {
    return (
        <div className="message">
            <div style={{ fontWeight: 'bold', marginBottom: '4px' }}>
                {senderName}
            </div>
            <div>{content}</div>
        </div>
    )
}
export default Message
