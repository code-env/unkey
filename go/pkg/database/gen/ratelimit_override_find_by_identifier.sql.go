// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: ratelimit_override_find_by_identifier.sql

package gen

import (
	"context"
)

const findRatelimitOverridesByIdentifier = `-- name: FindRatelimitOverridesByIdentifier :many
SELECT id, workspace_id, namespace_id, identifier, ` + "`" + `limit` + "`" + `, duration, async, sharding, created_at, updated_at, deleted_at FROM ratelimit_overrides
WHERE
    workspace_id = ?
    AND namespace_id = ?
    AND identifier LIKE ?
`

type FindRatelimitOverridesByIdentifierParams struct {
	WorkspaceID string `db:"workspace_id"`
	NamespaceID string `db:"namespace_id"`
	Identifier  string `db:"identifier"`
}

// FindRatelimitOverridesByIdentifier
//
//	SELECT id, workspace_id, namespace_id, identifier, `limit`, duration, async, sharding, created_at, updated_at, deleted_at FROM ratelimit_overrides
//	WHERE
//	    workspace_id = ?
//	    AND namespace_id = ?
//	    AND identifier LIKE ?
func (q *Queries) FindRatelimitOverridesByIdentifier(ctx context.Context, arg FindRatelimitOverridesByIdentifierParams) ([]RatelimitOverride, error) {
	rows, err := q.db.QueryContext(ctx, findRatelimitOverridesByIdentifier, arg.WorkspaceID, arg.NamespaceID, arg.Identifier)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []RatelimitOverride
	for rows.Next() {
		var i RatelimitOverride
		if err := rows.Scan(
			&i.ID,
			&i.WorkspaceID,
			&i.NamespaceID,
			&i.Identifier,
			&i.Limit,
			&i.Duration,
			&i.Async,
			&i.Sharding,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
