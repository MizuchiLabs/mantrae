-- name: ListAuditLogs :many
SELECT
  *
FROM
  audit_logs
WHERE
  profile_id = ?
ORDER BY
  created_at DESC
LIMIT
  ?
OFFSET
  ?;

-- name: CountAuditLogs :one
SELECT
  COUNT(*)
FROM
  audit_logs
WHERE
  profile_id = ?;

-- name: CreateAuditLog :exec
INSERT INTO
  audit_logs (
    profile_id,
    user_id,
    agent_id,
    event,
    details,
    created_at
  )
VALUES
  (?, ?, ?, ?, ?, CURRENT_TIMESTAMP);

-- name: DeleteOldAuditLogs :exec
DELETE FROM audit_logs
WHERE
  created_at < DATETIME ('now', '-90 days');
