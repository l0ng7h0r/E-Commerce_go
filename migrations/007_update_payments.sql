-- 007_update_payments.sql
ALTER TABLE payments ADD COLUMN IF NOT EXISTS method VARCHAR(50) DEFAULT 'phajay';
ALTER TABLE payments ADD COLUMN IF NOT EXISTS payment_url TEXT;

-- Add column logistic company, district to orders
ALTER TABLE orders ADD COLUMN IF NOT EXISTS logistic_company VARCHAR(50) DEFAULT '';
ALTER TABLE orders ADD COLUMN IF NOT EXISTS district VARCHAR(50) DEFAULT '';