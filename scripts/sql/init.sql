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

CREATE TABLE purchase_types (
    id   SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    cost INTEGER NOT NULL CHECK (cost >= 0)
);

CREATE TABLE IF NOT EXISTS purchases (    
    id               SERIAL PRIMARY KEY,
    purchaser_id     INTEGER NOT NULL,
    purchase_type_id INTEGER,
    FOREIGN KEY (purchase_type_id) REFERENCES purchase_types (id) ON DELETE SET NULL,
    FOREIGN KEY (purchaser_id) REFERENCES users (id) ON DELETE CASCADE
);

INSERT INTO purchase_types (name, cost) VALUES
    ('t-shirt', 80),
    ('cup', 20),
    ('book', 50),
    ('pen', 10),
    ('powerbank', 200),
    ('hoody', 300),
    ('umbrella', 200),
    ('socks', 10),
    ('wallet', 50),
    ('pink-hoody', 500)
ON CONFLICT (name) DO NOTHING;