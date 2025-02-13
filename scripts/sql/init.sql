CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    coins INTEGER NOT NULL DEFAULT 0 CHECK (coins >= 0),
    password_hash TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    sender_user_id INTEGER NOT NULL,
    receiver_user_id INTEGER NOT NULL,
    amount INT NOT NULL CHECK (amount > 0)
    FOREIGN KEY (sender_user_id) REFERENCES users (id) ON DELETE SET NULL
    FOREIGN KEY (receiver_user_id) REFERENCES users (id) ON DELETE SET NULL
);
