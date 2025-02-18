// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: permissions_by_key_id.sql

package gen

import (
	"context"
)

const findPermissionsForKey = `-- name: FindPermissionsForKey :many
WITH direct_permissions AS (
    SELECT p.name as permission_name
    FROM keys_permissions kp
    JOIN permissions p ON kp.permission_id = p.id
    WHERE kp.key_id = ?
),
role_permissions AS (
    SELECT p.name as permission_name
    FROM keys_roles kr
    JOIN roles_permissions rp ON kr.role_id = rp.role_id
    JOIN permissions p ON rp.permission_id = p.id
    WHERE kr.key_id = ?
)
SELECT DISTINCT permission_name
FROM (
    SELECT permission_name FROM direct_permissions
    UNION ALL
    SELECT permission_name FROM role_permissions
) all_permissions
`

type FindPermissionsForKeyParams struct {
	KeyID string `db:"key_id"`
}

// FindPermissionsForKey
//
//	WITH direct_permissions AS (
//	    SELECT p.name as permission_name
//	    FROM keys_permissions kp
//	    JOIN permissions p ON kp.permission_id = p.id
//	    WHERE kp.key_id = ?
//	),
//	role_permissions AS (
//	    SELECT p.name as permission_name
//	    FROM keys_roles kr
//	    JOIN roles_permissions rp ON kr.role_id = rp.role_id
//	    JOIN permissions p ON rp.permission_id = p.id
//	    WHERE kr.key_id = ?
//	)
//	SELECT DISTINCT permission_name
//	FROM (
//	    SELECT permission_name FROM direct_permissions
//	    UNION ALL
//	    SELECT permission_name FROM role_permissions
//	) all_permissions
func (q *Queries) FindPermissionsForKey(ctx context.Context, arg FindPermissionsForKeyParams) ([]string, error) {
	rows, err := q.db.QueryContext(ctx, findPermissionsForKey, arg.KeyID, arg.KeyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var permission_name string
		if err := rows.Scan(&permission_name); err != nil {
			return nil, err
		}
		items = append(items, permission_name)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
