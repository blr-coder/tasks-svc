ALTER TABLE tasks
    ADD COLUMN is_active BOOL NOT NULL DEFAULT true,
    ADD COLUMN currency  VARCHAR(10),
    ADD COLUMN amount     NUMERIC(20, 6)
