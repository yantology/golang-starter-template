-- Drop linkfy table and its trigger
DROP TRIGGER IF EXISTS update_linkfy_updated_at ON linkfy;
DROP TABLE IF EXISTS linkfy;