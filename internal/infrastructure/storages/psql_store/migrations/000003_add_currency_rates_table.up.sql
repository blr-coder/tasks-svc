CREATE TABLE IF NOT EXISTS currency_rates
(
    currency   VARCHAR(50) PRIMARY KEY,
    rate_eur   NUMERIC(20, 6) NOT NULL,
    created_at TIMESTAMP      NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    updated_at TIMESTAMP      NOT NULL DEFAULT (now() AT TIME ZONE 'utc')
);
