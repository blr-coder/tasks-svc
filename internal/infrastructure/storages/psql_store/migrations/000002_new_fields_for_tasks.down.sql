ALTER TABLE tasks
    DROP COLUMN IF EXISTS is_active,
    DROP COLUMN IF EXISTS currency,
    DROP COLUMN IF EXISTS amount
