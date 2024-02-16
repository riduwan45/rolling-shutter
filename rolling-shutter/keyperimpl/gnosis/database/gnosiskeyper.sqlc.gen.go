// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: gnosiskeyper.sql

package database

import (
	"context"

	"github.com/jackc/pgconn"
)

const getCurrentDecryptionTrigger = `-- name: GetCurrentDecryptionTrigger :one
SELECT eon, block, tx_pointer, identities_hash FROM current_decryption_trigger
WHERE eon = $1
`

func (q *Queries) GetCurrentDecryptionTrigger(ctx context.Context, eon int64) (CurrentDecryptionTrigger, error) {
	row := q.db.QueryRow(ctx, getCurrentDecryptionTrigger, eon)
	var i CurrentDecryptionTrigger
	err := row.Scan(
		&i.Eon,
		&i.Block,
		&i.TxPointer,
		&i.IdentitiesHash,
	)
	return i, err
}

const getSlotDecryptionSignatures = `-- name: GetSlotDecryptionSignatures :many
SELECT eon, block, keyper_index, tx_pointer, identities_hash, signature FROM slot_decryption_signatures
WHERE eon = $1 AND block = $2 AND tx_pointer = $3 AND identities_hash = $4
ORDER BY keyper_index ASC
LIMIT $5
`

type GetSlotDecryptionSignaturesParams struct {
	Eon            int64
	Block          int64
	TxPointer      int64
	IdentitiesHash []byte
	Limit          int32
}

func (q *Queries) GetSlotDecryptionSignatures(ctx context.Context, arg GetSlotDecryptionSignaturesParams) ([]SlotDecryptionSignature, error) {
	rows, err := q.db.Query(ctx, getSlotDecryptionSignatures,
		arg.Eon,
		arg.Block,
		arg.TxPointer,
		arg.IdentitiesHash,
		arg.Limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SlotDecryptionSignature
	for rows.Next() {
		var i SlotDecryptionSignature
		if err := rows.Scan(
			&i.Eon,
			&i.Block,
			&i.KeyperIndex,
			&i.TxPointer,
			&i.IdentitiesHash,
			&i.Signature,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTransactionSubmittedEventCount = `-- name: GetTransactionSubmittedEventCount :one
SELECT event_count FROM transaction_submitted_event_count
WHERE eon = $1
LIMIT 1
`

func (q *Queries) GetTransactionSubmittedEventCount(ctx context.Context, eon int64) (int64, error) {
	row := q.db.QueryRow(ctx, getTransactionSubmittedEventCount, eon)
	var event_count int64
	err := row.Scan(&event_count)
	return event_count, err
}

const getTransactionSubmittedEvents = `-- name: GetTransactionSubmittedEvents :many
SELECT index, block_number, block_hash, tx_index, log_index, eon, identity_prefix, sender, gas_limit FROM transaction_submitted_event
WHERE eon = $1 AND index >= $2
ORDER BY index ASC
LIMIT $3
`

type GetTransactionSubmittedEventsParams struct {
	Eon   int64
	Index int64
	Limit int32
}

func (q *Queries) GetTransactionSubmittedEvents(ctx context.Context, arg GetTransactionSubmittedEventsParams) ([]TransactionSubmittedEvent, error) {
	rows, err := q.db.Query(ctx, getTransactionSubmittedEvents, arg.Eon, arg.Index, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []TransactionSubmittedEvent
	for rows.Next() {
		var i TransactionSubmittedEvent
		if err := rows.Scan(
			&i.Index,
			&i.BlockNumber,
			&i.BlockHash,
			&i.TxIndex,
			&i.LogIndex,
			&i.Eon,
			&i.IdentityPrefix,
			&i.Sender,
			&i.GasLimit,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTransactionSubmittedEventsSyncedUntil = `-- name: GetTransactionSubmittedEventsSyncedUntil :one
SELECT block_number FROM transaction_submitted_events_synced_until LIMIT 1
`

func (q *Queries) GetTransactionSubmittedEventsSyncedUntil(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, getTransactionSubmittedEventsSyncedUntil)
	var block_number int64
	err := row.Scan(&block_number)
	return block_number, err
}

const getTxPointer = `-- name: GetTxPointer :one
SELECT eon, block, value FROM tx_pointer
WHERE eon = $1
`

func (q *Queries) GetTxPointer(ctx context.Context, eon int64) (TxPointer, error) {
	row := q.db.QueryRow(ctx, getTxPointer, eon)
	var i TxPointer
	err := row.Scan(&i.Eon, &i.Block, &i.Value)
	return i, err
}

const initTxPointer = `-- name: InitTxPointer :exec
INSERT INTO tx_pointer (eon, block, value)
VALUES ($1, $2, 0)
ON CONFLICT DO NOTHING
`

type InitTxPointerParams struct {
	Eon   int64
	Block int64
}

func (q *Queries) InitTxPointer(ctx context.Context, arg InitTxPointerParams) error {
	_, err := q.db.Exec(ctx, initTxPointer, arg.Eon, arg.Block)
	return err
}

const insertSlotDecryptionSignature = `-- name: InsertSlotDecryptionSignature :exec
INSERT INTO slot_decryption_signatures (eon, block, keyper_index, tx_pointer, identities_hash, signature)
VALUES ($1, $2, $3, $4, $5, $6)
`

type InsertSlotDecryptionSignatureParams struct {
	Eon            int64
	Block          int64
	KeyperIndex    int64
	TxPointer      int64
	IdentitiesHash []byte
	Signature      []byte
}

func (q *Queries) InsertSlotDecryptionSignature(ctx context.Context, arg InsertSlotDecryptionSignatureParams) error {
	_, err := q.db.Exec(ctx, insertSlotDecryptionSignature,
		arg.Eon,
		arg.Block,
		arg.KeyperIndex,
		arg.TxPointer,
		arg.IdentitiesHash,
		arg.Signature,
	)
	return err
}

const insertTransactionSubmittedEvent = `-- name: InsertTransactionSubmittedEvent :execresult
INSERT INTO transaction_submitted_event (
    index,
    block_number,
    block_hash,
    tx_index,
    log_index,
    eon,
    identity_prefix,
    sender,
    gas_limit
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
ON CONFLICT DO NOTHING
`

type InsertTransactionSubmittedEventParams struct {
	Index          int64
	BlockNumber    int64
	BlockHash      []byte
	TxIndex        int64
	LogIndex       int64
	Eon            int64
	IdentityPrefix []byte
	Sender         string
	GasLimit       int64
}

func (q *Queries) InsertTransactionSubmittedEvent(ctx context.Context, arg InsertTransactionSubmittedEventParams) (pgconn.CommandTag, error) {
	return q.db.Exec(ctx, insertTransactionSubmittedEvent,
		arg.Index,
		arg.BlockNumber,
		arg.BlockHash,
		arg.TxIndex,
		arg.LogIndex,
		arg.Eon,
		arg.IdentityPrefix,
		arg.Sender,
		arg.GasLimit,
	)
}

const setCurrentDecryptionTrigger = `-- name: SetCurrentDecryptionTrigger :exec
INSERT INTO current_decryption_trigger (eon, block, tx_pointer, identities_hash)
VALUES ($1, $2, $3, $4)
ON CONFLICT (eon) DO UPDATE
SET block = $2, tx_pointer = $3, identities_hash = $4
`

type SetCurrentDecryptionTriggerParams struct {
	Eon            int64
	Block          int64
	TxPointer      int64
	IdentitiesHash []byte
}

func (q *Queries) SetCurrentDecryptionTrigger(ctx context.Context, arg SetCurrentDecryptionTriggerParams) error {
	_, err := q.db.Exec(ctx, setCurrentDecryptionTrigger,
		arg.Eon,
		arg.Block,
		arg.TxPointer,
		arg.IdentitiesHash,
	)
	return err
}

const setTransactionSubmittedEventCount = `-- name: SetTransactionSubmittedEventCount :exec
INSERT INTO transaction_submitted_event_count (eon, event_count)
VALUES ($1, $2)
ON CONFLICT (eon) DO UPDATE
SET event_count = $2
`

type SetTransactionSubmittedEventCountParams struct {
	Eon        int64
	EventCount int64
}

func (q *Queries) SetTransactionSubmittedEventCount(ctx context.Context, arg SetTransactionSubmittedEventCountParams) error {
	_, err := q.db.Exec(ctx, setTransactionSubmittedEventCount, arg.Eon, arg.EventCount)
	return err
}

const setTransactionSubmittedEventsSyncedUntil = `-- name: SetTransactionSubmittedEventsSyncedUntil :exec
INSERT INTO transaction_submitted_events_synced_until (block_number) VALUES ($1)
ON CONFLICT (enforce_one_row) DO UPDATE
SET block_number = $1
`

func (q *Queries) SetTransactionSubmittedEventsSyncedUntil(ctx context.Context, blockNumber int64) error {
	_, err := q.db.Exec(ctx, setTransactionSubmittedEventsSyncedUntil, blockNumber)
	return err
}

const setTxPointer = `-- name: SetTxPointer :exec
INSERT INTO tx_pointer (eon, block, value)
VALUES ($1, $2, $3)
ON CONFLICT (eon) DO UPDATE
SET block = $2, value = $3
`

type SetTxPointerParams struct {
	Eon   int64
	Block int64
	Value int64
}

func (q *Queries) SetTxPointer(ctx context.Context, arg SetTxPointerParams) error {
	_, err := q.db.Exec(ctx, setTxPointer, arg.Eon, arg.Block, arg.Value)
	return err
}
