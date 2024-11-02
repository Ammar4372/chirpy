-- name: CreateUser :one
INSERT INTO users(Id, email, created_at, updated_at) VALUES (
    gen_random_uuid(),
    $1,
    $2,
    $3
) RETURNING *;

-- name: DeleteAllUser :exec
DELETE FROM users;