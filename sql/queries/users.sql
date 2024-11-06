-- name: CreateUser :one
INSERT INTO users(Id, email, hashed_password, created_at, updated_at) VALUES (
    gen_random_uuid(),
    $1,
    $2,
    $3,
    $4
) RETURNING Id, email, created_at, updated_at;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: DeleteAllUser :exec
DELETE FROM users;