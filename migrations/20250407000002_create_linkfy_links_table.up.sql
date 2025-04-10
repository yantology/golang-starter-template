-- Create linkfy_links table
CREATE TABLE linkfy_links (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    linkfy_id UUID NOT NULL REFERENCES linkfy(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    name_url VARCHAR(255) NOT NULL,
    icons_url TEXT,
    display_order INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create index on linkfy_id for faster lookups
CREATE INDEX idx_linkfy_links_linkfy_id ON linkfy_links(linkfy_id);