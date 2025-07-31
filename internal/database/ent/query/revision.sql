-- sqlfluff:dialect:postgres
-- sqlfluff:max_line_length:1024
-- sqlfluff:rules:capitalisation.keywords:capitalisation_policy:upper

SET search_path TO public;

-- Creates a table named 'aurora_schema_revisions' with the following columns:
-- name: CreateTableRevisions :exec
CREATE TABLE IF NOT EXISTS aurora_schema_revisions (
    -- primary key column
    id TEXT PRIMARY KEY,
    -- aurora_schema_revisions name column
    description TEXT NOT NULL,
    -- execution timestamp column
    executed_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    -- execution time column
    execution_time BIGINT NOT NULL DEFAULT 0
);
-- Retrieves a row from the table 'aurora_schema_revisions' with option ':one'
-- name: GetRevision :one
SELECT
    id,
    description,
    executed_at,
    execution_time
FROM
    aurora_schema_revisions
WHERE
    id = sqlc.arg(id);

-- Inserts a row into the table 'aurora_schema_revisions' with option ':one'
-- name: InsertRevision :one
INSERT INTO aurora_schema_revisions (
    id,
    description,
    executed_at,
    execution_time
) VALUES (
    sqlc.arg(id),
    sqlc.arg(description),
    sqlc.arg(executed_at),
    sqlc.arg(execution_time)
)
RETURNING *;

-- Inserts a row into the table 'aurora_schema_revisions' with option ':exec'
-- name: ExecInsertRevision :exec
INSERT INTO aurora_schema_revisions (
    id,
    description,
    executed_at,
    execution_time
) VALUES (
    sqlc.arg(id),
    sqlc.arg(description),
    sqlc.arg(executed_at),
    sqlc.arg(execution_time)
);

-- Retrieves a list of rows from the table 'aurora_schema_revisions' with option ':many'
-- name: ListRevisions :many
SELECT
    id,
    description,
    executed_at,
    execution_time
FROM
    aurora_schema_revisions
ORDER BY
    id
LIMIT
    sqlc.narg(page_limit)::INT
    OFFSET
    sqlc.narg(page_offset)::INT;
