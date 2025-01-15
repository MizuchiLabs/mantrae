-- name: CreateTraefikConfig :exec
INSERT INTO
    traefik (profile_id, source, entrypoints, overview, config)
VALUES
    (?, ?, ?, ?, ?);

-- name: GetTraefikConfig :one
SELECT
    *
FROM
    traefik
WHERE
    id = ?;

-- name: GetTraefikConfigBySource :one
SELECT
    *
FROM
    traefik
WHERE
    profile_id = ?
    AND source = ?;

-- name: UpdateTraefikConfig :exec
UPDATE traefik
SET
    source = ?,
    entrypoints = ?,
    overview = ?,
    config = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = ?;

-- name: DeleteTraefikConfig :exec
DELETE FROM traefik
WHERE
    id = ?;
