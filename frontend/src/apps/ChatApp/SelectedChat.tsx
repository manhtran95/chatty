import MessageList from './MessageList';

interface MessageData {
  senderName: string;
  content: string;
}

interface SelectedChatProps {
  chatId: string | null;
  messages: MessageData[];
}

function SelectedChat({ chatId, messages }: SelectedChatProps) {
  if (chatId === null) {
    return (
      <div style={{ width: '75%', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
        <div style={{ color: '#666' }}>Select a chat to view messages</div>
      </div>
    );
  }

  return (
    <div style={{ width: '75%', height: '100vh' }}>
      <MessageList messages={messages} />
    </div>
  );
}

export default SelectedChat;