-- name: GetNote :one
SELECT *
FROM note
WHERE id = $1;

-- name: CreateNote :one
INSERT INTO note (
    content, created_at, expires_at, expires_at_timezone
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: DeleteNote :exec
DELETE FROM note
WHERE id = $1;
