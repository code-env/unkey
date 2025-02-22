// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: key_find_by_hash.sql

package gen

import (
	"context"
)

const findKeyByHash = `-- name: FindKeyByHash :one
SELECT
    id, key_auth_id, hash, start, workspace_id, for_workspace_id, name, owner_id, identity_id, meta, created_at, expires, created_at_m, updated_at_m, deleted_at_m, deleted_at, refill_day, refill_amount, last_refill_at, enabled, remaining_requests, ratelimit_async, ratelimit_limit, ratelimit_duration, environment
FROM ` + "`" + `keys` + "`" + `
WHERE hash = ?
`

// FindKeyByHash
//
//	SELECT
//	    id, key_auth_id, hash, start, workspace_id, for_workspace_id, name, owner_id, identity_id, meta, created_at, expires, created_at_m, updated_at_m, deleted_at_m, deleted_at, refill_day, refill_amount, last_refill_at, enabled, remaining_requests, ratelimit_async, ratelimit_limit, ratelimit_duration, environment
//	FROM `keys`
//	WHERE hash = ?
func (q *Queries) FindKeyByHash(ctx context.Context, hash string) (Key, error) {
	row := q.db.QueryRowContext(ctx, findKeyByHash, hash)
	var i Key
	err := row.Scan(
		&i.ID,
		&i.KeyAuthID,
		&i.Hash,
		&i.Start,
		&i.WorkspaceID,
		&i.ForWorkspaceID,
		&i.Name,
		&i.OwnerID,
		&i.IdentityID,
		&i.Meta,
		&i.CreatedAt,
		&i.Expires,
		&i.CreatedAtM,
		&i.UpdatedAtM,
		&i.DeletedAtM,
		&i.DeletedAt,
		&i.RefillDay,
		&i.RefillAmount,
		&i.LastRefillAt,
		&i.Enabled,
		&i.RemainingRequests,
		&i.RatelimitAsync,
		&i.RatelimitLimit,
		&i.RatelimitDuration,
		&i.Environment,
	)
	return i, err
}
