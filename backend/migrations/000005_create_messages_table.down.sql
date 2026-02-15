-- Drop the trigger first (depends on function and table)
DROP TRIGGER IF EXISTS trigger_update_chat_updated_at ON messages;

-- Drop the trigger function (depends on nothing)
DROP FUNCTION IF EXISTS update_chat_updated_at();

-- Drop the index (depends on table)
DROP INDEX IF EXISTS idx_messages_chat_created;

-- Drop the messages table (last, as other things depend on it)
DROP TABLE IF EXISTS messages;