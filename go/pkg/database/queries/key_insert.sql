-- name: InsertKey :exec
INSERT INTO `keys` (
    id,
    key_auth_id,
    hash,
    start,
    workspace_id,
    for_workspace_id,
    name,
    owner_id,
    identity_id,
    meta,
    created_at,
    expires,
    created_at_m,
    enabled,
    remaining_requests,
    ratelimit_async,
    ratelimit_limit,
    ratelimit_duration,
    environment
) VALUES (
    sqlc.arg(id),
    sqlc.arg(keyring_id),
    sqlc.arg(hash),
    sqlc.arg(start),
    sqlc.arg(workspace_id),
    sqlc.arg(for_workspace_id),
    sqlc.arg(name),
    null,
    sqlc.arg(identity_id),
    sqlc.arg(meta),
    sqlc.arg(created_at),
    sqlc.arg(expires),
    UNIX_TIMESTAMP() * 1000,
    sqlc.arg(enabled),
    sqlc.arg(remaining_requests),
    sqlc.arg(ratelimit_async),
    sqlc.arg(ratelimit_limit),
    sqlc.arg(ratelimit_duration),
    sqlc.arg(environment)
);
