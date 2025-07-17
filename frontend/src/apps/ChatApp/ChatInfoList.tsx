import ChatInfo from './ChatInfo';

interface ChatData {
  id: string;
  name: string;
  lastMessage: string;
}

interface ChatInfoListProps {
  chats: ChatData[];
  selectedChatId: string | null;
  onChatSelect: (chatId: string) => void;
}

function ChatInfoList({ chats, selectedChatId, onChatSelect }: ChatInfoListProps) {
  return (
    <div style={{ width: '25%', borderRight: '1px solid #ddd', height: '100vh', overflowY: 'auto' }}>
      {chats.map((chat) => (
        <ChatInfo
          key={chat.id}
          chatId={chat.id}
          name={chat.name}
          lastMessage={chat.lastMessage}
          isSelected={selectedChatId === chat.id}
          onClick={() => onChatSelect(chat.id)}
        />
      ))}
    </div>
  );
}

export default ChatInfoList;