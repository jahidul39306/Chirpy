-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(),
    NOW(), -- created_at
    NOW(), -- updated_at
    $1,  -- email
    $2 -- password
)
RETURNING *;

-- name: DeleteUsers :exec
DELETE FROM users;