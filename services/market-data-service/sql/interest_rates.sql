CREATE TABLE IF NOT EXISTS interest_rates (
    id SERIAL PRIMARY KEY,
    series_id VARCHAR(20) NOT NULL,
    series_name TEXT NOT NULL,
    date DATE NOT NULL,
    rate NUMERIC(8, 4),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(series_id, date)
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_interest_rates_series_date ON interest_rates(series_id, date);
CREATE INDEX IF NOT EXISTS idx_interest_rates_date ON interest_rates(date);

-- Common FRED series IDs:
-- FEDFUNDS: Federal Funds Rate
-- DGS10: 10-Year Treasury Constant Maturity Rate
-- DGS2: 2-Year Treasury Constant Maturity Rate
-- DGS1: 1-Year Treasury Constant Maturity Rate
-- TB3MS: 3-Month Treasury Bill Rate 