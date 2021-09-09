// Code generated by sqlc. DO NOT EDIT.
// source: query.sql

package kprdb

import (
	"context"
)

const getDecryptionKey = `-- name: GetDecryptionKey :one
SELECT epoch_id, keyper_index, decryption_key FROM keyper.decryption_key
WHERE epoch_id = $1
`

func (q *Queries) GetDecryptionKey(ctx context.Context, epochID int64) (KeyperDecryptionKey, error) {
	row := q.db.QueryRowContext(ctx, getDecryptionKey, epochID)
	var i KeyperDecryptionKey
	err := row.Scan(&i.EpochID, &i.KeyperIndex, &i.DecryptionKey)
	return i, err
}