CREATE TABLE
    messages (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        sender_id UUID NOT NULL,
        chat_id UUID NOT NULL,
        content TEXT NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (sender_id) REFERENCES users (id) ON DELETE CASCADE,
        FOREIGN KEY (chat_id) REFERENCES chats (id) ON DELETE CASCADE
    );

-- Create composite index for efficient message retrieval by chat and time
CREATE INDEX idx_messages_chat_created ON messages (chat_id, created_at, id);


-- Create trigger function to update chat's updated_at when a new message is inserted
CREATE OR REPLACE FUNCTION update_chat_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE chats
    SET updated_at = GREATEST(updated_at, NEW.created_at)
    WHERE id = NEW.chat_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger that fires after each message insert
CREATE TRIGGER trigger_update_chat_updated_at
    AFTER INSERT ON messages
    FOR EACH ROW
    EXECUTE FUNCTION update_chat_updated_at();