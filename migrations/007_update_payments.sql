-- 007_update_payments.sql
ALTER TABLE payments ADD COLUMN IF NOT EXISTS method VARCHAR(50) DEFAULT 'phajay';
ALTER TABLE payments ADD COLUMN IF NOT EXISTS payment_url TEXT;
