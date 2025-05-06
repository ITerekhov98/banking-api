CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    from_account_id INTEGER REFERENCES accounts(id),
    to_account_id INTEGER REFERENCES accounts(id),
    amount NUMERIC(20, 2) NOT NULL CHECK (amount > 0),
    created_at TIMESTAMP DEFAULT now()
);