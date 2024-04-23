-- name: GetNote :one
SELECT *
FROM note
WHERE id = $1;

-- name: CreateNote :one
INSERT INTO note (
    content, created_at, expires_at, expires_at_timezone, key_hash
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: DeleteNote :exec
DELETE FROM note
WHERE id = $1;
