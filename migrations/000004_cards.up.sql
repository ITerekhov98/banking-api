CREATE TABLE IF NOT EXISTS cards (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    account_id INTEGER REFERENCES accounts(id),
    number BYTEA NOT NULL,      -- PGP-шифрованный номер
    expiry BYTEA NOT NULL,      -- PGP-шифрованная дата
    cvv_hash TEXT NOT NULL,     -- bcrypt хеш CVV
    hmac TEXT NOT NULL,         -- HMAC всей карты
    created_at TIMESTAMP DEFAULT now()
);