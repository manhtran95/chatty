CREATE TABLE
    messages (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        sender_id UUID NOT NULL,
        chat_id UUID NOT NULL,
        content TEXT NOT NULL,
        created_at TIMESTAMP
        WITH
            TIME ZONE DEFAULT NOW (),
            FOREIGN KEY (sender_id) REFERENCES users (id) ON DELETE CASCADE,
            FOREIGN KEY (chat_id) REFERENCES chats (id) ON DELETE CASCADE
    );

-- Create composite index for efficient message retrieval by chat and time
CREATE INDEX idx_messages_chat_created ON messages (chat_id, created_at);