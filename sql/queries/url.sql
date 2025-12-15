-- name: SaveURL :one
INSERT INTO url (url, alias, expire_at)
VALUES ($1, $2, $3)
RETURNING id
;

-- name: GetURL :one
SELECT url FROM url
WHERE alias = $1
AND (expire_at IS NULL OR expire_at > now())
;

-- name: DeleteURL :execrows
DELETE from url
WHERE alias = $1
;