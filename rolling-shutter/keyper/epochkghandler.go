package keyper

import (
	"context"
	"log"

	"github.com/pkg/errors"

	"github.com/shutter-network/shutter/shlib/puredkg"
	"github.com/shutter-network/shutter/shuttermint/keyper/epochkg"
	"github.com/shutter-network/shutter/shuttermint/keyper/kprdb"
	"github.com/shutter-network/shutter/shuttermint/shdb"
	"github.com/shutter-network/shutter/shuttermint/shmsg"
)

type epochKGHandler struct {
	config Config
	db     *kprdb.Queries
}

func (h *epochKGHandler) handleDecryptionTrigger(ctx context.Context, msg *decryptionTrigger) ([]shmsg.P2PMessage, error) {
	eon, err := h.db.GetEonForEpoch(ctx, shdb.EncodeUint64(msg.EpochID))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get eon for epoch %d from db", msg.EpochID)
	}
	batchConfig, err := h.db.GetBatchConfig(ctx, int32(eon.ConfigIndex))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get config %d from db", eon.ConfigIndex)
	}

	// get our keyper index (and check that we in fact are a keyper)
	encodedAddress := shdb.EncodeAddress(h.config.Address())
	keyperIndex := int64(-1)
	for i, address := range batchConfig.Keypers {
		if address == encodedAddress {
			keyperIndex = int64(i)
			break
		}
	}
	if keyperIndex == -1 {
		log.Printf("ignoring decryption trigger for epoch %d as we are not a keyper", msg.EpochID)
		return nil, nil
	}

	// check if we already computed (and therefore most likely sent) our key share
	shareExists, err := h.db.ExistsDecryptionKeyShare(ctx, kprdb.ExistsDecryptionKeyShareParams{
		EpochID:     shdb.EncodeUint64(msg.EpochID),
		KeyperIndex: keyperIndex,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get decryption key share for epoch %d from db", msg.EpochID)
	}
	if shareExists {
		return nil, nil // we already sent our share
	}

	// fetch dkg result from db
	dkgResultDB, err := h.db.GetDKGResult(ctx, eon.Eon)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get dkg result for eon %d from db", eon.Eon)
	}
	if !dkgResultDB.Success {
		log.Printf("ignoring decryption trigger for epoch %d as eon key generation failed", msg.EpochID)
		return nil, nil
	}
	pureDKGResult, err := shdb.DecodePureDKGResult(dkgResultDB.PureResult)
	if err != nil {
		return nil, err
	}

	// compute the key share
	epochKG := epochkg.NewEpochKG(pureDKGResult)
	share := epochKG.ComputeEpochSecretKeyShare(msg.EpochID)
	encodedShare := share.Marshal()
	if err != nil {
		return nil, err
	}

	// store share in db and sent it
	err = h.db.InsertDecryptionKeyShare(ctx, kprdb.InsertDecryptionKeyShareParams{
		EpochID:            shdb.EncodeUint64(msg.EpochID),
		KeyperIndex:        keyperIndex,
		DecryptionKeyShare: encodedShare,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to insert decryption key share")
	}
	log.Printf("sending decryption key share for epoch %d", msg.EpochID)
	return []shmsg.P2PMessage{
		&shmsg.DecryptionKeyShare{
			InstanceID:  h.config.InstanceID,
			EpochID:     msg.EpochID,
			KeyperIndex: uint64(keyperIndex),
			Share:       encodedShare,
		},
	}, nil
}

func (h *epochKGHandler) insertDecryptionKeyShare(ctx context.Context, msg *decryptionKeyShare) error {
	encodedShare := shdb.EncodeEpochSecretKeyShare(msg.share)
	err := h.db.InsertDecryptionKeyShare(ctx, kprdb.InsertDecryptionKeyShareParams{
		EpochID:            shdb.EncodeUint64(msg.epochID),
		KeyperIndex:        int64(msg.keyperIndex),
		DecryptionKeyShare: encodedShare,
	})
	if err != nil {
		return errors.Wrapf(
			err,
			"failed to insert decryption key share for epoch %d from keyper %d",
			msg.epochID,
			msg.keyperIndex,
		)
	}
	return nil
}

func (h *epochKGHandler) aggregateDecryptionKeySharesFromDB(
	ctx context.Context,
	pureDKGResult *puredkg.Result,
	epochID uint64,
) (*epochkg.EpochKG, error) {
	shares, err := h.db.SelectDecryptionKeyShares(ctx, shdb.EncodeUint64(epochID))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get decryption key shares for epoch %d from db", epochID)
	}

	epochKG := epochkg.NewEpochKG(pureDKGResult)
	// For simplicity, we aggregate shares even if we don't have enough of them yet.
	for _, share := range shares {
		shareDecoded, err := shdb.DecodeEpochSecretKeyShare(share.DecryptionKeyShare)
		if err != nil {
			log.Printf(
				"Warning: invalid decryption key share in db for epoch %d and keyper %d",
				share.EpochID,
				share.KeyperIndex,
			)
			continue
		}
		err = epochKG.HandleEpochSecretKeyShare(&epochkg.EpochSecretKeyShare{
			Eon:    pureDKGResult.Eon,
			Epoch:  epochID,
			Sender: uint64(share.KeyperIndex),
			Share:  shareDecoded,
		})
		if err != nil {
			log.Printf(
				"error processing decryption key share for epoch %d of keyper %d: %s",
				shdb.DecodeUint64(share.EpochID), share.KeyperIndex, err,
			)
			continue
		}
	}

	return epochKG, nil
}

func (h *epochKGHandler) handleDecryptionKeyShare(ctx context.Context, msg *decryptionKeyShare) ([]shmsg.P2PMessage, error) {
	// Insert the share into the db. We assume that it's valid as it already passed the libp2p
	// validator.
	if err := h.insertDecryptionKeyShare(ctx, msg); err != nil {
		return nil, err
	}

	// Check that we don't know the decryption key yet
	keyExists, err := h.db.ExistsDecryptionKey(ctx, shdb.EncodeUint64(msg.epochID))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to query decryption key for epoch %d", msg.epochID)
	}
	if keyExists {
		return nil, nil
	}

	// fetch dkg result from db
	eon, err := h.db.GetEonForEpoch(ctx, shdb.EncodeUint64(msg.epochID))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get eon for epoch %d from db", msg.epochID)
	}
	dkgResultDB, err := h.db.GetDKGResult(ctx, eon.Eon)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get dkg result for eon %d from db", eon.Eon)
	}
	if !dkgResultDB.Success {
		log.Printf("ignoring decryption trigger for epoch %d as eon key generation failed", msg.epochID)
		return nil, nil
	}
	pureDKGResult, err := shdb.DecodePureDKGResult(dkgResultDB.PureResult)
	if err != nil {
		return nil, err
	}

	// aggregate epoch secret key
	epochKG, err := h.aggregateDecryptionKeySharesFromDB(ctx, pureDKGResult, msg.epochID)
	if err != nil {
		return nil, err
	}
	decryptionKey, ok := epochKG.SecretKeys[msg.epochID]
	if !ok {
		numShares := uint64(len(epochKG.SecretShares))
		if numShares < pureDKGResult.Threshold {
			// not enough shares yet
			return nil, nil
		}
		return nil, errors.Errorf(
			"failed to generate decryption key for epoch %d even though we have enough shares",
			msg.epochID,
		)
	}
	decryptionKeyEncoded := decryptionKey.Marshal()

	// send decryption key
	tag, err := h.db.InsertDecryptionKey(ctx, kprdb.InsertDecryptionKeyParams{
		EpochID:       shdb.EncodeUint64(msg.epochID),
		DecryptionKey: decryptionKeyEncoded,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to store decryption key for epoch %d in db", msg.epochID)
	}
	if tag.RowsAffected() == 0 {
		log.Printf("attempted to insert decryption key for epoch %d, but it already exists", msg.epochID)
		return nil, nil
	}
	log.Printf("broadcasting decryption key for epoch %d", msg.epochID)
	return []shmsg.P2PMessage{
		&shmsg.DecryptionKey{
			InstanceID: h.config.InstanceID,
			EpochID:    msg.epochID,
			Key:        decryptionKeyEncoded,
		},
	}, nil
}

func (h *epochKGHandler) handleDecryptionKey(ctx context.Context, msg *decryptionKey) ([]shmsg.P2PMessage, error) {
	// Insert the key into the db. We assume that it's valid as it already passed the libp2p
	// validator.
	encodedKey := msg.key.Marshal()
	tag, err := h.db.InsertDecryptionKey(ctx, kprdb.InsertDecryptionKeyParams{
		EpochID:       shdb.EncodeUint64(msg.epochID),
		DecryptionKey: encodedKey,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to insert decryption key for epoch %d", msg.epochID)
	}
	if tag.RowsAffected() == 0 {
		log.Printf("attempted to insert decryption key for epoch %d, but it already exists", msg.epochID)
	}
	return nil, nil
}
