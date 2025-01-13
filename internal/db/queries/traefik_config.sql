-- name: CreateTraefikConfig :exec
INSERT INTO
    traefik_config (
        profile_id,
        source,
        entrypoints,
        overview,
        config,
        last_sync
    )
VALUES
    (?, ?, ?, ?, ?, ?);

-- name: GetTraefikConfig :one
SELECT
    *
FROM
    traefik_config
WHERE
    id = ?;

-- name: GetTraefikConfigBySource :one
SELECT
    *
FROM
    traefik_config
WHERE
    profile_id = ?
    AND source = ?;

-- name: UpdateTraefikConfig :exec
UPDATE traefik_config
SET
    source = ?,
    entrypoints = ?,
    overview = ?,
    config = ?,
    last_sync = ?
WHERE
    id = ?;

-- name: DeleteTraefikConfig :exec
DELETE FROM traefik_config
WHERE
    id = ?;

-- name: GetTraefikConfigLastSync :one
SELECT
    last_sync
FROM
    traefik_config
WHERE
    id = ?;

-- name: UpdateTraefikConfigLastSync :exec
UPDATE traefik_config
SET
    last_sync = CURRENT_TIMESTAMP
WHERE
    id = ?;
