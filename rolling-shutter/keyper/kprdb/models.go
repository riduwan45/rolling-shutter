// Code generated by sqlc. DO NOT EDIT.

package kprdb

import (
	"database/sql"
	"time"
)

type KeyperDecryptionKey struct {
	EpochID       []byte
	KeyperIndex   sql.NullInt64
	DecryptionKey []byte
}

type KeyperDecryptionKeyShare struct {
	EpochID            []byte
	KeyperIndex        int64
	DecryptionKeyShare []byte
}

type KeyperDecryptionTrigger struct {
	EpochID []byte
}

type KeyperEon struct {
	Eon         int64
	Height      int64
	BatchIndex  []byte
	ConfigIndex int64
}

type KeyperMetaInf struct {
	Key   string
	Value string
}

type KeyperPolyEval struct {
	Eon             int64
	ReceiverAddress string
	Eval            []byte
}

type KeyperPuredkg struct {
	Eon     int64
	Puredkg []byte
}

type KeyperTendermintBatchConfig struct {
	ConfigIndex int32
	Height      int64
	Keypers     []string
	Threshold   int32
}

type KeyperTendermintEncryptionKey struct {
	Address             string
	EncryptionPublicKey []byte
}

type KeyperTendermintOutgoingMessage struct {
	ID          int32
	Description string
	Msg         []byte
}

type KeyperTendermintSyncMetum struct {
	CurrentBlock        int64
	LastCommittedHeight int64
	SyncTimestamp       time.Time
}
