CREATE TABLE users (
    id UUID PRIMARY KEY,
    name VARCHAR(100),
    password_hash VARCHAR(64) NOT NULL,
    salt VARCHAR(32) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE TABLE wallets (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    balance NUMERIC(20,2) DEFAULT 0,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE TABLE transactions (
    id UUID PRIMARY KEY,
    wallet_id UUID REFERENCES wallets(id),
    amount NUMERIC(20,2),
    type VARCHAR(50),
    description TEXT,
    related_user_id UUID,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);
