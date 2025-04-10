-- Create linkfy table
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE linkfy (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    username VARCHAR(255) UNIQUE NOT NULL,
    avatar_url TEXT,
    name VARCHAR(255) NOT NULL,
    bio TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create trigger on linkfy table for updated_at
CREATE TRIGGER update_linkfy_updated_at
    BEFORE UPDATE ON linkfy
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();