-- Таблица для учета стримов
CREATE TABLE IF NOT EXISTS streams (
    id VARCHAR(50) PRIMARY KEY,       -- Twitch Stream ID
    channel VARCHAR(50) NOT NULL,      -- Имя канала (например, 'silvername')
    started_at TIMESTAMP WITH TIME ZONE NOT NULL,
    is_active BOOLEAN DEFAULT TRUE     -- Активен ли стрим сейчас
);

-- Таблица для 5-минутной статистики слов
CREATE TABLE IF NOT EXISTS word_stats (
    id SERIAL PRIMARY KEY,
    stream_id VARCHAR(50) REFERENCES streams(id) ON DELETE CASCADE,
    time_frame TIMESTAMP WITH TIME ZONE NOT NULL, -- Время начала 5-минутки
    word VARCHAR(100) NOT NULL,
    count INT NOT NULL,
    
    -- Уникальный индекс, чтобы при повторной вставке мы могли обновить количество (UPSERT)
    CONSTRAINT unique_stream_time_word UNIQUE (stream_id, time_frame, word)
);

-- Индекс для быстрой выборки статистики конкретного стрима по времени
CREATE INDEX IF NOT EXISTS idx_word_stats_stream_time ON word_stats(stream_id, time_frame);