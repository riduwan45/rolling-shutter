// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0

package database

import ()

type TransactionSubmittedEvent struct {
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

type TransactionSubmittedEventCount struct {
	Eon        int64
	EventCount int64
}

type TransactionSubmittedEventsSyncedUntil struct {
	EnforceOneRow bool
	BlockNumber   int64
}

type TxPointer struct {
	Eon   int64
	Block int64
	Value int64
}
