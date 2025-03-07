ALTER TABLE users
DROP COLUMN is_admin,
DROP COLUMN usage_limit_secs,
DROP COLUMN usage_total_secs,
DROP COLUMN password_hash;

-- Удалить админа с telegram_id = 0
DELETE FROM users WHERE telegram_id = 0;