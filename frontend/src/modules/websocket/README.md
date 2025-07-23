# WebSocket Implementation

This directory contains the WebSocket implementation for the chat application.

## Architecture

### WebSocketService (`../../services/WebSocketService.ts`)

- Core WebSocket connection management
- Handles connection, disconnection, and message sending
- Provides message handler registration system
- Connection status tracking

### WebSocketContext (`WebSocketContext.tsx`)

- React context for WebSocket state management
- Auto-connects when user is authenticated
- Provides WebSocket functionality to components
- Handles cleanup on unmount

### useWebSocketMessage (`useWebSocketMessage.ts`)

- Custom hook for easier WebSocket message handling
- Provides connection status and message callbacks

## Usage

### Basic Usage in Components

```tsx
import { useWebSocket } from '../../modules/websocket/WebSocketContext'

function MyComponent() {
    const { isConnected, sendMessage, addMessageHandler } = useWebSocket()

    useEffect(() => {
        const removeHandler = addMessageHandler((message) => {
            console.log('Received:', message)
        })

        return () => removeHandler()
    }, [addMessageHandler])

    const handleSendMessage = () => {
        sendMessage({
            type: 'chat_message',
            data: { content: 'Hello!' },
        })
    }

    return (
        <div>
            <p>Status: {isConnected ? 'Connected' : 'Disconnected'}</p>
            <button onClick={handleSendMessage}>Send Message</button>
        </div>
    )
}
```

### Using the Custom Hook

```tsx
import { useWebSocketMessage } from '../../modules/websocket/useWebSocketMessage'

function MyComponent() {
    const { isConnected } = useWebSocketMessage({
        onMessage: (message) => {
            console.log('Received message:', message)
        },
        onConnect: () => {
            console.log('WebSocket connected')
        },
        onDisconnect: () => {
            console.log('WebSocket disconnected')
        },
    })

    return <div>Connected: {isConnected}</div>
}
```

## Connection Management

### Automatic Connection

- WebSocket automatically connects when user is authenticated
- Connection is established using the user's access token
- Connection status is tracked and displayed in the UI

### Graceful Disconnection

- Connection is automatically closed when user logs out
- Cleanup is performed on component unmount
- Connection status is properly reset

### Error Handling

- Connection errors are logged and handled gracefully
- Failed connections are retried automatically
- Message sending is protected against disconnected state

## Message Format

Messages sent and received should follow this format:

```typescript
interface WebSocketMessage {
    type: string // Message type (e.g., 'new_message', 'chat_update')
    data: any // Message payload
}
```

## Security

- WebSocket connection requires authentication token
- Token is automatically included in connection URL
- Connection is closed when user logs out

## Best Practices

1. **Always check connection status** before sending messages
2. **Use the cleanup function** returned by `addMessageHandler`
3. **Handle connection errors** gracefully
4. **Use the custom hook** for simpler message handling
5. **Test connection status** in your components
