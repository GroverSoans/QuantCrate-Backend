CREATE TABLE IF NOT EXISTS fgi_index (
    id SERIAL PRIMARY KEY,
    value INT NOT NULL,                        -- E.g. 47
    value_classification TEXT NOT NULL,        -- "Neutral", "Fear", etc.
    timestamp_unix BIGINT NOT NULL,            -- Raw UNIX timestamp
    resolved_date DATE NOT NULL,               -- Converted date (UTC)
    UNIQUE (resolved_date)                     -- One entry per day
);
