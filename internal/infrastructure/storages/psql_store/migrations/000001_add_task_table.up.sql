CREATE TYPE task_status_t AS ENUM ('PENDING', 'PROCESSING', 'DONE');
COMMENT ON TYPE task_status_t IS e'Статусы задачи:'
    '\n * PENDING    - в ожидани'
    '\n * PROCESSING - в процессе обработки'
    '\n * DONE       - обработка завершена';

CREATE TABLE IF NOT EXISTS tasks
(
    id          BIGSERIAL PRIMARY KEY,
    title       VARCHAR(150) NOT NULL UNIQUE,
    description TEXT,
    customer_id UUID          NOT NULL,
    executor_id UUID,
    status      task_status_t NOT NULL,
    created_at  TIMESTAMP     NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    updated_at  TIMESTAMP     NOT NULL DEFAULT (now() AT TIME ZONE 'utc')
);
