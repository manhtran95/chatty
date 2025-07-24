import { messageClasses } from '../../utils/tailwindClasses'

interface MessageProps {
    senderName: string
    content: string
    isOwnMessage?: boolean
    timestamp?: string
}

function Message({ senderName, content, isOwnMessage = false, timestamp }: MessageProps) {
    // TODO: IMPLEMENT THIS
    isOwnMessage = content.length % 2 === 0
    return (
        <div className={isOwnMessage ? messageClasses.ownMessage : messageClasses.container}>
            <div className={messageClasses.sender}>
                {senderName}
            </div>
            <div className={messageClasses.content}>
                {content}
            </div>
            {timestamp && (
                <div className={messageClasses.timestamp}>
                    {timestamp}
                </div>
            )}
        </div>
    )
}

export default Message
