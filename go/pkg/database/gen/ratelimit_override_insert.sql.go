// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: ratelimit_override_insert.sql

package gen

import (
	"context"
	"time"
)

const insertRatelimitOverride = `-- name: InsertRatelimitOverride :exec
INSERT INTO
    ` + "`" + `ratelimit_overrides` + "`" + ` (
        id,
        workspace_id,
        namespace_id,
        identifier,
        ` + "`" + `limit` + "`" + `,
        duration,
        async,
        created_at
    )
VALUES
    (
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        false,
         ?
    )
`

type InsertRatelimitOverrideParams struct {
	ID          string    `db:"id"`
	WorkspaceID string    `db:"workspace_id"`
	NamespaceID string    `db:"namespace_id"`
	Identifier  string    `db:"identifier"`
	Limit       int32     `db:"limit"`
	Duration    int32     `db:"duration"`
	CreatedAt   time.Time `db:"created_at"`
}

// InsertRatelimitOverride
//
//	INSERT INTO
//	    `ratelimit_overrides` (
//	        id,
//	        workspace_id,
//	        namespace_id,
//	        identifier,
//	        `limit`,
//	        duration,
//	        async,
//	        created_at
//	    )
//	VALUES
//	    (
//	        ?,
//	        ?,
//	        ?,
//	        ?,
//	        ?,
//	        ?,
//	        false,
//	         ?
//	    )
func (q *Queries) InsertRatelimitOverride(ctx context.Context, arg InsertRatelimitOverrideParams) error {
	_, err := q.db.ExecContext(ctx, insertRatelimitOverride,
		arg.ID,
		arg.WorkspaceID,
		arg.NamespaceID,
		arg.Identifier,
		arg.Limit,
		arg.Duration,
		arg.CreatedAt,
	)
	return err
}
