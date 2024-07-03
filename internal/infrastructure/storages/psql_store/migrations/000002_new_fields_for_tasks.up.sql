ALTER TABLE tasks
    ADD COLUMN is_active BOOL NOT NULL DEFAULT true,
    ADD COLUMN currency  VARCHAR(10),
    ADD COLUMN price     NUMERIC(20, 6)
