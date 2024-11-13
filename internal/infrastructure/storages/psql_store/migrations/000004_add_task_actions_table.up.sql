CREATE TYPE action_type_t AS ENUM ('COMMIT', 'MERGE_REQUEST', 'MERGE');
COMMENT ON TYPE task_status_t IS e'Action types:'
    '\n * COMMIT        - commit'
    '\n * MERGE_REQUEST - merge request'
    '\n * MERGE         - merge';

CREATE TABLE IF NOT EXISTS task_actions
(
    id          BIGSERIAL PRIMARY KEY,
    external_id BIGSERIAL     NOT NULL UNIQUE,
    task_id     BIGSERIAL     NOT NULL REFERENCES tasks (id) ON UPDATE RESTRICT ON DELETE RESTRICT,
    type        action_type_t NOT NULL,
    url         TEXT          NOT NULL,
    created_at  TIMESTAMP     NOT NULL DEFAULT (now() AT TIME ZONE 'utc')
);
