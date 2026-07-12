-- name: CreateUser :one
INSERT INTO users (
    id,
    created_at,
    updated_at,
    email,
    hashed_password
)
VALUES (
    gen_random_uuid(),
    now(),
    now(),
    $1,
    $2
)
RETURNING (
    id,
    created_at,
    updated_at,
    email
);


-- name: DeleteUsers :exec
DELETE FROM users;


-- name: GetUserByEmail :one
SELECT
    id,
    email,
    hashed_password
FROM
    users
WHERE
    email = $1;
