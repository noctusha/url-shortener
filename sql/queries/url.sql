-- name: SaveURL :one
INSERT INTO url (url, alias)
VALUES ($1, $2)
RETURNING id
;

-- name: GetURL :one
SELECT url FROM url
WHERE alias = $1
;

-- name: DeleteURL :exec
DELETE from url
WHERE alias = $1
;