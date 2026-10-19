-- 009_add_name_phone_shop_name_to_users.sql
-- Add name, phone, and shop_name columns to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS name VARCHAR(255) DEFAULT '';
ALTER TABLE users ADD COLUMN IF NOT EXISTS phone VARCHAR(50) DEFAULT '';
ALTER TABLE users ADD COLUMN IF NOT EXISTS shop_name VARCHAR(255) DEFAULT '';
