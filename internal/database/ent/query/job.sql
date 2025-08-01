-- sqlfluff:dialect:postgres
-- sqlfluff:max_line_length:1024
-- sqlfluff:rules:capitalisation.keywords:capitalisation_policy:upper

SET search_path TO public;

-- Waits for a job to complete by its ID.
-- name: WaitForJob :one
SELECT sys.wait_for_job(sqlc.arg(job_id)::TEXT)::BOOLEAN AS ok;
