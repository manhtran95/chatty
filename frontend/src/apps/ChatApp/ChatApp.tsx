import { useState, useEffect } from 'react'
import ChatInfoList from './ChatInfoList'
import SelectedChat from './SelectedChat'
import type { ChatData, NewMessage } from './types'
import { receiveNewMessage } from './utils/commandHandlers'

function ChatApp() {
    // states
    const [selectedChatId, setSelectedChatId] = useState<string | null>(null)
    const [chatListData, setChatListData] = useState<ChatData[] | null>(null)

    // stub data
    const stubInitData = [
        {
            "id": "1",
            "name": "Chat0",
            "userList": [
                {
                    "id": "3a092631-cbcc-4ab1-a1ae-d617da0050d4",
                    "username": "carol50"
                },
                {
                    "id": "693fe40f-3a17-4689-b51b-da057980b990",
                    "username": "eve68"
                }
            ],
            "lastMessage": "See you soon.",
            "messages": [
                {
                    "senderName": "eve68",
                    "content": "See you soon."
                },
                {
                    "senderName": "eve68",
                    "content": "Check this out."
                },
                {
                    "senderName": "carol50",
                    "content": "Any updates?"
                },
                {
                    "senderName": "eve68",
                    "content": "Hey!"
                }
            ]
        },
        {
            "id": "2",
            "name": "Chat1",
            "userList": [
                {
                    "id": "c78c1108-600c-4c99-8ec7-187e0b7f0150",
                    "username": "alice86"
                },
                {
                    "id": "950bdece-46c7-4b08-9e10-573bf0d4382e",
                    "username": "dave91"
                }
            ],
            "lastMessage": "How's it going?",
            "messages": [
                {
                    "senderName": "alice86",
                    "content": "How's it going?"
                },
                {
                    "senderName": "dave91",
                    "content": "Thanks!"
                },
                {
                    "senderName": "dave91",
                    "content": "Let's meet tomorrow."
                }
            ]
        },
        {
            "id": "3",
            "name": "Chat2",
            "userList": [
                {
                    "id": "05ec6324-a3b8-476d-a474-6617fcf00d63",
                    "username": "eve32"
                },
                {
                    "id": "57d856fe-7996-4a3f-814d-38cc91ac3c13",
                    "username": "bob86"
                }
            ],
            "lastMessage": "What's up?",
            "messages": [
                {
                    "senderName": "eve32",
                    "content": "What's up?"
                },
                {
                    "senderName": "bob86",
                    "content": "Thanks!"
                },
                {
                    "senderName": "eve32",
                    "content": "Good morning!"
                }
            ]
        }
    ]

    console.log(`selectedChatId: ${selectedChatId}`)

    // run once when the component mounts
    useEffect(() => {
        setChatListData(stubInitData)
        // Optional cleanup
        return () => {
            console.log('Component unmounted');
        };
    }, []); // Empty dependency array = run once

    // deduced data
    const chatInfoList = chatListData ? chatListData.map((chat) => ({
        id: chat.id,
        name: chat.name,
        lastMessage: chat.lastMessage || 'No messages'
    })) : []
    const selectedChatMessages: Array<{
        senderName: string;
        content: string;
    }> = selectedChatId !== null && chatListData ? chatListData.find(chat => chat.id === selectedChatId)?.messages || [] : []

    return (
        <>
            <div style={{ display: 'flex', height: '80vh' }}>
                <ChatInfoList
                    chats={chatInfoList}
                    selectedChatId={selectedChatId}
                    onChatSelect={setSelectedChatId}
                />
                <SelectedChat
                    chatId={selectedChatId}
                    messages={selectedChatMessages}
                />
            </div>
            <div>
                <button onClick={() => receiveNewMessage({
                    "chatId": "1",
                    "senderName": "manh",
                    "content": `New message here!${Math.random()}`
                }, setChatListData)}>Receive new mess testing</button>
            </div>
        </>
    )
}

export default ChatApp
