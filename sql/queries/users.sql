-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
    )
    RETURNING *;
 
-- name: DeleteAllUsers :exec
DELETE FROM users;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: UpdateUserAuthentication :one
UPDATE users
SET email = $1, hashed_password = $2, updated_at = NOW()
WHERE id = $3
RETURNING *;

-- name: UpdateUserIsChirpyRedTrue :exec
UPDATE users
SET is_chirpy_red = TRUE
WHERE id = $1;

-- name: UpdateUserIsChirpyRedFalse :exec
UPDATE users
SET is_chirpy_red = FALSE
WHERE id = $1;
