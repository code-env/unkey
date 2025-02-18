// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: ratelimit_namespace_delete.sql

package gen

import (
	"context"
	"database/sql"
)

const deleteRatelimitNamespace = `-- name: DeleteRatelimitNamespace :execresult
UPDATE ` + "`" + `ratelimit_namespaces` + "`" + `
SET deleted_at = ?
WHERE id = ?
`

type DeleteRatelimitNamespaceParams struct {
	Now sql.NullTime `db:"now"`
	ID  string       `db:"id"`
}

// DeleteRatelimitNamespace
//
//	UPDATE `ratelimit_namespaces`
//	SET deleted_at = ?
//	WHERE id = ?
func (q *Queries) DeleteRatelimitNamespace(ctx context.Context, arg DeleteRatelimitNamespaceParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, deleteRatelimitNamespace, arg.Now, arg.ID)
}
