-- sqlfluff:dialect:postgres
-- sqlfluff:max_line_length:1024
-- sqlfluff:rules:capitalisation.keywords:capitalisation_policy:upper

SET search_path TO public;

-- Creates a table named 'aurora_schema_revisions' with the following columns:
-- name: CreateTableRevisions :exec
CREATE TABLE IF NOT EXISTS aurora_schema_revisions (
    -- primary key column
    id TEXT PRIMARY KEY,
    -- revision name
    description TEXT NOT NULL,
    -- total number of statements
    total INT NOT NULL DEFAULT 0,
    -- count of statements executed in this revision
    count INT NOT NULL DEFAULT 0,
    -- error issued during the execution of the revision
    error TEXT NULL,
    -- error_stmt is the statement that caused the error
    error_stmt TEXT NULL,
    -- execution timestamp column
    executed_at TIMESTAMP WITH TIME ZONE NOT NULL,
    -- execution time column
    execution_time BIGINT NOT NULL DEFAULT 0
);

-- Retrieves a row from the table 'aurora_schema_revisions' with option ':one'
-- name: GetRevision :one
SELECT
    id,
    description,
    total,
    count,
    error,
    error_stmt,
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
    total,
    count,
    error,
    error_stmt,
    executed_at,
    execution_time
) VALUES (
    sqlc.arg(id),
    sqlc.arg(description),
    sqlc.arg(total),
    sqlc.arg(count),
    sqlc.narg(error),
    sqlc.narg(error_stmt),
    sqlc.arg(executed_at),
    sqlc.arg(execution_time)
)
RETURNING *;

-- Inserts a row into the table 'aurora_schema_revisions' with option ':exec'
-- name: ExecInsertRevision :exec
INSERT INTO aurora_schema_revisions (
    id,
    description,
    total,
    count,
    error,
    error_stmt,
    executed_at,
    execution_time
) VALUES (
    sqlc.arg(id),
    sqlc.arg(description),
    sqlc.arg(total),
    sqlc.arg(count),
    sqlc.narg(error),
    sqlc.narg(error_stmt),
    sqlc.arg(executed_at),
    sqlc.arg(execution_time)
);

-- Upserts a row into the table 'aurora_schema_revisions' with option ':one'
-- name: UpsertRevision :one
INSERT INTO aurora_schema_revisions (
    id,
    description,
    total,
    count,
    error,
    error_stmt,
    executed_at,
    execution_time
) VALUES (
    sqlc.arg(id),
    sqlc.arg(description),
    sqlc.arg(total),
    sqlc.arg(count),
    sqlc.narg(error),
    sqlc.narg(error_stmt),
    sqlc.arg(executed_at),
    sqlc.arg(execution_time)
)
ON CONFLICT (id) DO UPDATE SET id = sqlc.arg(id)
RETURNING *;

-- Upserts a row into the table 'aurora_schema_revisions' with option ':exec'
-- name: ExecUpsertRevision :exec
INSERT INTO aurora_schema_revisions (
    id,
    description,
    total,
    count,
    error,
    error_stmt,
    executed_at,
    execution_time
) VALUES (
    sqlc.arg(id),
    sqlc.arg(description),
    sqlc.arg(total),
    sqlc.arg(count),
    sqlc.narg(error),
    sqlc.narg(error_stmt),
    sqlc.arg(executed_at),
    sqlc.arg(execution_time)
)
ON CONFLICT (id) DO UPDATE SET id = sqlc.arg(id);


-- Updates a row in the table 'revision' with option ':one'
-- name: UpdateRevision :one
UPDATE aurora_schema_revisions
SET
    description = CASE
        WHEN 'description' = ANY(sqlc.arg(update_mask)::TEXT [])
            THEN sqlc.arg(description)
        ELSE description
    END,
    total = CASE
        WHEN 'total' = ANY(sqlc.arg(update_mask)::TEXT [])
            THEN sqlc.arg(total)
        ELSE total
    END,
    count = CASE
        WHEN 'count' = ANY(sqlc.arg(update_mask)::TEXT [])
            THEN sqlc.arg(count)
        ELSE count
    END,
    error = CASE
        WHEN 'error' = ANY(sqlc.arg(update_mask)::TEXT [])
            THEN sqlc.narg(error)
        ELSE error
    END,
    error_stmt = CASE
        WHEN 'error_stmt' = ANY(sqlc.arg(update_mask)::TEXT [])
            THEN sqlc.narg(error_stmt)
        ELSE error_stmt
    END,
    executed_at = CASE
        WHEN 'executed_at' = ANY(sqlc.arg(update_mask)::TEXT [])
            THEN sqlc.arg(executed_at)
        ELSE executed_at
    END,
    execution_time = CASE
        WHEN 'execution_time' = ANY(sqlc.arg(update_mask)::TEXT [])
            THEN sqlc.arg(execution_time)
        ELSE execution_time
    END
WHERE
    id = sqlc.arg(id)
RETURNING *;

-- Updates a row in the table 'revision' with option ':exec'
-- name: ExecUpdateRevision :exec
UPDATE aurora_schema_revisions
SET
    description = CASE
        WHEN 'description' = ANY(sqlc.arg(update_mask)::TEXT [])
            THEN sqlc.arg(description)
        ELSE description
    END,
    total = CASE
        WHEN 'total' = ANY(sqlc.arg(update_mask)::TEXT [])
            THEN sqlc.arg(total)
        ELSE total
    END,
    count = CASE
        WHEN 'count' = ANY(sqlc.arg(update_mask)::TEXT [])
            THEN sqlc.arg(count)
        ELSE count
    END,
    error = CASE
        WHEN 'error' = ANY(sqlc.arg(update_mask)::TEXT [])
            THEN sqlc.narg(error)
        ELSE error
    END,
    error_stmt = CASE
        WHEN 'error_stmt' = ANY(sqlc.arg(update_mask)::TEXT [])
            THEN sqlc.narg(error_stmt)
        ELSE error_stmt
    END,
    executed_at = CASE
        WHEN 'executed_at' = ANY(sqlc.arg(update_mask)::TEXT [])
            THEN sqlc.arg(executed_at)
        ELSE executed_at
    END,
    execution_time = CASE
        WHEN 'execution_time' = ANY(sqlc.arg(update_mask)::TEXT [])
            THEN sqlc.arg(execution_time)
        ELSE execution_time
    END
WHERE
    id = sqlc.arg(id);

-- Deletes a row from the table 'aurora_schema_revisions' with option ':one'
-- name: DeleteRevision :one
DELETE FROM aurora_schema_revisions
WHERE id = sqlc.arg(id)
RETURNING *;

-- Deletes a row from the table 'aurora_schema_revisions' with option ':exec'
-- name: ExecDeleteRevision :exec
DELETE FROM aurora_schema_revisions
WHERE id = sqlc.arg(id);

-- Retrieves a list of rows from the table 'aurora_schema_revisions' with option ':many'
-- name: ListRevisions :many
SELECT
    id,
    description,
    total,
    count,
    error,
    error_stmt,
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
