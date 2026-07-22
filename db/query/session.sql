
-- name: CreateSession :one
INSERT INTO sessions (
    id,
    username,
    user_agent ,
    user_ip ,
    refresh_token,
    is_blocked,
    refresh_token_expired_at
) VALUES (
  $1, $2,$3,$4,$5,$6,$7
)
RETURNING *;

-- name: GetSession :one
SELECT * FROM sessions
WHERE id = $1 LIMIT 1;