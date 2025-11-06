-- name: CreateJob :one
INSERT INTO jobs (
    id,
    type,
    payload,
    status,
    max_attempts
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetJob :one
SELECT * FROM jobs
WHERE id = $1;

-- name: DequeueJob :one
UPDATE jobs
SET 
    status = 'processing',
    updated_at = NOW()
WHERE id = (
    SELECT id
    FROM jobs
    WHERE status = 'pending'
        AND scheduled_at <= NOW()
    ORDER BY scheduled_at ASC
    LIMIT 1
    FOR UPDATE SKIP LOCKED
)
RETURNING *;

-- name: CompleteJob :exec
UPDATE jobs
SET 
    status = 'completed',
    updated_at = NOW()
WHERE id = $1;

-- name: FailJob :exec
UPDATE jobs
SET 
    status = $2,
    attempts = $3,
    error_message = $4,
    scheduled_at = $5,
    updated_at = NOW()
WHERE id = $1;

-- name: ListJobs :many
SELECT * FROM jobs
WHERE status = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountJobsByStatus :one
SELECT COUNT(*) FROM jobs
WHERE status = $1;