CREATE TABLE IF NOT EXISTS users (
                                       id BIGSERIAL PRIMARY KEY,
                                       telegram_chat_id BIGINT NOT NULL UNIQUE,
                                       timezone TEXT NOT NULL DEFAULT 'UTC',
                                       created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS reminders (
                                      id BIGSERIAL PRIMARY KEY,
                                      user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                      text TEXT NOT NULL,
                                      notify_at TIMESTAMPTZ NOT NULL,
                                      status VARCHAR(20) NOT NULL DEFAULT 'pending',
                                      created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_reminders_status_notify_at ON reminders (status, notify_at);