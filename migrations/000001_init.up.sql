CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'student',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS events (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    date TIMESTAMP NOT NULL,
    location VARCHAR(255),
    max_participants INTEGER DEFAULT 0,
    creator_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS event_registrations (
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    event_id INTEGER REFERENCES events(id) ON DELETE CASCADE,
    registered_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, event_id)
);

-- Seed initial data
-- Admin account (password: admin123)
INSERT INTO users (id, email, password_hash, role) VALUES 
(1, 'admin@goevent.ru', '$2a$10$LyOz/xKXxd/vFurYAsMj7.nrl5PD.SNeMUnAMMm2oDV2SiGwoTVFq', 'admin'),
(2, 'student@goevent.ru', '$2a$10$LyOz/xKXxd/vFurYAsMj7.nrl5PD.SNeMUnAMMm2oDV2SiGwoTVFq', 'student')
ON CONFLICT (email) DO NOTHING;

-- Reset standard serial sequences so automatic ids don't conflict
SELECT setval('users_id_seq', (SELECT MAX(id) FROM users));

-- Sample events
INSERT INTO events (id, title, description, date, location, max_participants, creator_id) VALUES
(1, 'Хакатон GoEvent 2026', 'Ежегодный студенческий хакатон по разработке высоконагруженных систем на Go. Ценные призы, мерч и стажировки для лучших участников!', '2026-06-15 10:00:00', 'Коворкинг IT-кластера, зал А', 50, 1),
(2, 'Лекция: Архитектура микросервисов', 'Приглашенный спикер из крупной IT-компании расскажет про паттерны проектирования микросервисов, распределенные транзакции и кэширование.', '2026-06-20 14:00:00', 'Аудитория 404, главный корпус', 100, 1),
(3, 'Спортивный турнир по настольному теннису', 'Студенческий турнир для всех желающих. Приходите поддержать друзей или побороться за кубок университета!', '2026-06-25 16:30:00', 'Спортзал №2', 32, 1)
ON CONFLICT (id) DO NOTHING;

SELECT setval('events_id_seq', (SELECT MAX(id) FROM events));

