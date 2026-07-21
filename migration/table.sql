CREATE TABLE IF NOT EXISTS streams (
    id VARCHAR(50) PRIMARY KEY,       -- Twitch Stream ID
    channel VARCHAR(50) NOT NULL,      -- Имя канала (например, 'silvername')
    started_at TIMESTAMP WITH TIME ZONE NOT NULL,
    is_active BOOLEAN DEFAULT TRUE     -- Активен ли стрим сейчас
);


CREATE TABLE IF NOT EXISTS word_stats (
    id SERIAL PRIMARY KEY,
    stream_id VARCHAR(50) REFERENCES streams(id) ON DELETE CASCADE,
    time_frame TIMESTAMP WITH TIME ZONE NOT NULL,
    word VARCHAR(100) NOT NULL,
    count INT NOT NULL,
    CONSTRAINT unique_stream_time_word UNIQUE (stream_id, time_frame, word)
);

CREATE TABLE user_period_stats (
    id SERIAL PRIMARY KEY,
    stream_id VARCHAR(255) NOT NULL,
    username VARCHAR(255) NOT NULL,
    msg_count INT NOT NULL DEFAULT 0,
    period_start TIMESTAMP NOT NULL,
    UNIQUE(stream_id, username)
);


CREATE INDEX IF NOT EXISTS idx_word_stats_stream_time ON word_stats(stream_id, time_frame);
