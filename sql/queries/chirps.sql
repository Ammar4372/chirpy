-- name: CreateChrip :one
INSERT INTO chirps (Id, body, user_id, created_at, updated_at) VALUES (
    gen_random_uuid(),
    $1,
    $2,
    $3,
    $4
) RETURNING *;

-- name: GetAllChirps :many
SELECT * FROM chirps ORDER BY created_at;

-- name: GetChirpById :one
SELECT * FROM chirps WHERE Id = $1 ORDER BY created_at;

-- name: GetChirpByAuthorId :many
SELECT * FROM chirps WHERE user_id = $1 ORDER BY created_at;
