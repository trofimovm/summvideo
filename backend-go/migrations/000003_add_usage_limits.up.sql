ALTER TABLE users
ADD COLUMN is_admin BOOLEAN DEFAULT FALSE,
ADD COLUMN usage_limit_secs INTEGER DEFAULT 3600, -- 1 час (3600 секунд) по умолчанию
ADD COLUMN usage_total_secs INTEGER DEFAULT 0,
ADD COLUMN password_hash VARCHAR(255) DEFAULT NULL;

-- Создаем админа с заданными кредами
INSERT INTO users (
    telegram_id, 
    username, 
    first_name, 
    last_name, 
    is_admin, 
    is_active,
    password_hash
) VALUES (
    0, -- специальный ID для админа
    'superadmin', 
    'Super', 
    'Admin', 
    TRUE, 
    TRUE,
    '$2a$10$HcuzvB7hFXB.IUvNmXZQyeyyPpBbhFh3/rC5F6eZeEKqm.g9kUjdK' -- хеш для пароля fRecuRioNihi
);