-- name: ListAuditLogs :many
SELECT
  a.id,
  a.profile_id,
  p.name AS profile_name,
  a.user_id,
  u.username AS user_name,
  a.agent_id,
  ag.hostname AS agent_name,
  a.event,
  a.details,
  a.created_at
FROM
  audit_logs a
  LEFT JOIN profiles p ON a.profile_id = p.id
  LEFT JOIN users u ON a.user_id = u.id
  LEFT JOIN agents ag ON a.agent_id = ag.id
ORDER BY
  a.created_at DESC
LIMIT
  ?
OFFSET
  ?;

-- name: CountAuditLogs :one
SELECT
  COUNT(*)
FROM
  audit_logs;

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
