import Message from './Message';

interface MessageData {
  senderName: string;
  content: string;
}

interface MessageListProps {
  messages: MessageData[];
}

function MessageList({ messages }: MessageListProps) {
  return (
    <div style={{ height: '100%', overflowY: 'auto', paddingLeft: '16px', paddingRight: '16px', display: 'flex', flexDirection: 'column-reverse' }}>
      {messages.map((message, index) => (
        <Message
          key={index}
          senderName={message.senderName}
          content={message.content}
        />
      ))}
    </div>
  );
}

export default MessageList;