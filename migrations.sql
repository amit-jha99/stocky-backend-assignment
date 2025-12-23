CREATE TABLE reward_events (
    id UUID PRIMARY KEY,
    user_id VARCHAR(50) NOT NULL,
    stock_symbol VARCHAR(20) NOT NULL,
    quantity NUMERIC(18,6) NOT NULL,
    rewarded_at TIMESTAMP NOT NULL,
    idempotency_key VARCHAR(100) UNIQUE
);

CREATE TABLE ledger_entries (
    id UUID PRIMARY KEY,
    event_id UUID REFERENCES reward_events(id),
    entry_type VARCHAR(20),
    symbol VARCHAR(20),
    amount NUMERIC(18,4),
    direction VARCHAR(10)
);

-- Cached stock prices
CREATE TABLE stock_prices (
    stock_symbol VARCHAR(20) PRIMARY KEY,
    price_inr NUMERIC(18,4) NOT NULL,
    fetched_at TIMESTAMP NOT NULL
);
