-- sqlfluff:dialect:postgres
-- sqlfluff:max_line_length:1024
-- sqlfluff:rules:capitalisation.keywords:capitalisation_policy:upper

SET search_path TO public;

-- Creates a table named 'aurora_schema_locks' with the following columns:
-- name: CreateTableLocks :exec
CREATE TABLE IF NOT EXISTS aurora_schema_locks (
    -- primary key column
    id TEXT PRIMARY KEY,
    -- execution timestamp column
    created_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- Retrieves a row from the table 'aurora_schema_locks' with option ':one'
-- name: GetLock :one
SELECT
    id,
    created_at
FROM
    aurora_schema_locks
WHERE
    id = sqlc.arg(id);

-- Inserts a row into the table 'aurora_schema_locks' with option ':one'
-- name: InsertLock :one
INSERT INTO aurora_schema_locks (
    id,
    created_at
) VALUES (
    sqlc.arg(id),
    sqlc.arg(created_at)
)
RETURNING *;

-- Inserts a row into the table 'aurora_schema_locks' with option ':exec'
-- name: ExecInsertLock :exec
INSERT INTO aurora_schema_locks (
    id,
    created_at
) VALUES (
    sqlc.arg(id),
    sqlc.arg(created_at)
);

-- Deletes a row from the table 'aurora_schema_locks' with option ':one'
-- name: DeleteLock :one
DELETE FROM aurora_schema_locks
WHERE id = sqlc.arg(id)
RETURNING *;

-- Deletes a row from the table 'aurora_schema_locks' with option ':exec'
-- name: ExecDeleteLock :exec
DELETE FROM aurora_schema_locks
WHERE id = sqlc.arg(id);
