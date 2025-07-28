-- sqlfluff:dialect:postgres
-- sqlfluff:max_line_length:1024
-- sqlfluff:rules:capitalisation.keywords:capitalisation_policy:upper

SET search_path TO public;

-- The schema 'sys' is created to hold system-related tables.
-- name: CreateSchemaSys :exec
CREATE SCHEMA IF NOT EXISTS sys;

-- Creates a table named 'sys.jobs' with the following columns:
-- The table 'sys.jobs' is created to track jobs in the system.
-- name: CreateTableJobs :exec
CREATE TABLE IF NOT EXISTS sys.jobs (
    -- primary key column
    job_id TEXT PRIMARY KEY,
    -- revision name
    status TEXT NOT NULL,
    -- total number of statements
    details TEXT NULL
);

-- Inserts a row into the table 'sys.jobs' with option ':one'
-- name: InsertJob :one
INSERT INTO sys.jobs (
    job_id,
    status,
    details
) VALUES (
    sqlc.arg(job_id),
    sqlc.arg(status),
    sqlc.narg(details)
)
RETURNING *;

-- Inserts a row into the table 'sys.jobs' with option ':exec'
-- name: ExecInsertJob :exec
INSERT INTO sys.jobs (
    job_id,
    status,
    details
) VALUES (
    sqlc.arg(job_id),
    sqlc.arg(status),
    sqlc.narg(details)
);

-- Retrieves a row from the table 'sys.jobs' with option ':one'
-- name: GetJob :one
SELECT
    job_id,
    status,
    details
FROM
    sys.jobs
WHERE
    job_id = sqlc.arg(job_id);

-- Deletes a row from the table 'sys.jobs' with option ':one'
-- name: DeleteJob :one
DELETE FROM sys.jobs
WHERE job_id = sqlc.arg(job_id)
RETURNING *;

-- Deletes a row from the table 'sys.jobs' with option ':exec'
-- name: ExecDeleteJob :exec
DELETE FROM sys.jobs
WHERE job_id = sqlc.arg(job_id);
