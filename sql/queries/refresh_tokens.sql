-- name: CreateRefreshToken :exec
INSERT INTO refresh_tokens (
    token, created_at, updated_at, user_id, expires_at, revoked_at
)
VALUES (
    $1,
    NOW(),
    NOW(),
    $2,
    NOW() + INTERVAL '60 days',
    NULL
);


-- name: GetUserFromRefreshToken :one
SELECT
    user_id,
    expires_at,
    revoked_at
FROM refresh_tokens
WHERE $1 = token;
