// Code generated by sqlc. DO NOT EDIT.
// source: query.sql

package cltrdb

import (
	"context"
)

const getBiggestEpochID = `-- name: GetBiggestEpochID :one
SELECT epoch_id FROM collator.epoch_id ORDER BY epoch_id DESC LIMIT 1
`

func (q *Queries) GetBiggestEpochID(ctx context.Context) ([]byte, error) {
	row := q.db.QueryRow(ctx, getBiggestEpochID)
	var epoch_id []byte
	err := row.Scan(&epoch_id)
	return epoch_id, err
}

const getLastBatchEpochID = `-- name: GetLastBatchEpochID :one
SELECT epoch_id FROM collator.decryption_trigger ORDER BY epoch_id DESC LIMIT 1
`

func (q *Queries) GetLastBatchEpochID(ctx context.Context) ([]byte, error) {
	row := q.db.QueryRow(ctx, getLastBatchEpochID)
	var epoch_id []byte
	err := row.Scan(&epoch_id)
	return epoch_id, err
}

const getMeta = `-- name: GetMeta :one
SELECT key, value FROM collator.meta_inf WHERE key = $1
`

func (q *Queries) GetMeta(ctx context.Context, key string) (CollatorMetaInf, error) {
	row := q.db.QueryRow(ctx, getMeta, key)
	var i CollatorMetaInf
	err := row.Scan(&i.Key, &i.Value)
	return i, err
}

const getTransactionsByEpoch = `-- name: GetTransactionsByEpoch :many
SELECT encrypted_tx FROM collator.transaction WHERE epoch_id = $1 ORDER BY tx_id
`

func (q *Queries) GetTransactionsByEpoch(ctx context.Context, epochID []byte) ([][]byte, error) {
	rows, err := q.db.Query(ctx, getTransactionsByEpoch, epochID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items [][]byte
	for rows.Next() {
		var encrypted_tx []byte
		if err := rows.Scan(&encrypted_tx); err != nil {
			return nil, err
		}
		items = append(items, encrypted_tx)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTrigger = `-- name: GetTrigger :one
SELECT epoch_id, batch_hash FROM collator.decryption_trigger WHERE epoch_id = $1
`

func (q *Queries) GetTrigger(ctx context.Context, epochID []byte) (CollatorDecryptionTrigger, error) {
	row := q.db.QueryRow(ctx, getTrigger, epochID)
	var i CollatorDecryptionTrigger
	err := row.Scan(&i.EpochID, &i.BatchHash)
	return i, err
}

const insertEpochID = `-- name: InsertEpochID :exec
INSERT INTO collator.epoch_id (epoch_id) VALUES ($1)
`

func (q *Queries) InsertEpochID(ctx context.Context, epochID []byte) error {
	_, err := q.db.Exec(ctx, insertEpochID, epochID)
	return err
}

const insertMeta = `-- name: InsertMeta :exec
INSERT INTO collator.meta_inf (key, value) VALUES ($1, $2)
`

type InsertMetaParams struct {
	Key   string
	Value string
}

func (q *Queries) InsertMeta(ctx context.Context, arg InsertMetaParams) error {
	_, err := q.db.Exec(ctx, insertMeta, arg.Key, arg.Value)
	return err
}

const insertTrigger = `-- name: InsertTrigger :exec
INSERT INTO collator.decryption_trigger (epoch_id, batch_hash) VALUES ($1, $2)
`

type InsertTriggerParams struct {
	EpochID   []byte
	BatchHash []byte
}

func (q *Queries) InsertTrigger(ctx context.Context, arg InsertTriggerParams) error {
	_, err := q.db.Exec(ctx, insertTrigger, arg.EpochID, arg.BatchHash)
	return err
}

const insertTx = `-- name: InsertTx :exec
INSERT INTO collator.transaction (tx_id, epoch_id, encrypted_tx) VALUES ($1, $2, $3)
`

type InsertTxParams struct {
	TxID        []byte
	EpochID     []byte
	EncryptedTx []byte
}

func (q *Queries) InsertTx(ctx context.Context, arg InsertTxParams) error {
	_, err := q.db.Exec(ctx, insertTx, arg.TxID, arg.EpochID, arg.EncryptedTx)
	return err
}
