CREATE TABLE IF NOT EXISTS credits (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    account_id INTEGER NOT NULL REFERENCES accounts(id),
    principal NUMERIC(20,2) NOT NULL,
    interest_rate NUMERIC(5,2) NOT NULL,
    term_months INTEGER NOT NULL,
    monthly_payment NUMERIC(20,2) NOT NULL,
    created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE IF NOT EXISTS payment_schedules (
    id SERIAL PRIMARY KEY,
    credit_id INTEGER REFERENCES credits(id) ON DELETE CASCADE,
    due_date DATE NOT NULL,
    amount NUMERIC(20,2) NOT NULL,
    paid BOOLEAN DEFAULT false,
    penalty_applied BOOLEAN DEFAULT false
);