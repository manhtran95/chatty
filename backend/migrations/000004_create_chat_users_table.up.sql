CREATE TABLE
    chat_users (
        chat_id UUID NOT NULL,
        user_id UUID NOT NULL,
        PRIMARY KEY (chat_id, user_id),
        FOREIGN KEY (chat_id) REFERENCES chats (id) ON DELETE CASCADE,
        FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
    );

-- Create hash indexes for better performance on UUID lookups
CREATE INDEX idx_chat_users_chat_id ON chat_users USING hash (chat_id);

CREATE INDEX idx_chat_users_user_id ON chat_users USING hash (user_id);