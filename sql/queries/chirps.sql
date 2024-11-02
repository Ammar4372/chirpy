-- name: CreateChrip :one
INSERT INTO chirps (Id, body, user_id, created_at, updated_at) VALUES (
    gen_random_uuid(),
    $1,
    $2,
    $3,
    $4
) RETURNING *;

-- name: GetAllChirps :many
SELECT * FROM chirps;

-- name: GetChirpById :one
SELECT * FROM chirps WHERE Id = $1;