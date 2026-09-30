-- 008_add_qr_code_to_payments.sql
-- Adds qr_code column to store the raw EMVCo QR payload for in-app QR scanning
ALTER TABLE payments ADD COLUMN IF NOT EXISTS qr_code TEXT;
