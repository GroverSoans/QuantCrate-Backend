CREATE TABLE IF NOT EXISTS ohlcv (
    id SERIAL PRIMARY KEY,
    ticker VARCHAR(10) NOT NULL,
    date DATE NOT NULL,
    open NUMERIC(12, 6),
    high NUMERIC(12, 6),
    low NUMERIC(12, 6),
    close NUMERIC(12, 6),
    volume NUMERIC,
    UNIQUE (ticker, date)
);

CREATE INDEX idx_ticker_date ON ohlcv(ticker, date);