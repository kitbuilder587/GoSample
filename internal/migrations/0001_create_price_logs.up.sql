-- migrations/0001_create_price_logs.up.sql
CREATE TABLE IF NOT EXISTS price_logs (
                                          id SERIAL PRIMARY KEY,
                                          coin TEXT NOT NULL,
                                          price_usd NUMERIC(20, 10) NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL DEFAULT now()
    );
