package mocksequencer

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	txtypes "github.com/shutter-network/txtypes/types"
)

type blockData struct {
	baseFee  *big.Int
	gasLimit uint64
}

type SequencerProcessor struct {
	port       int16
	nonces     map[string]map[string]uint64
	collators  map[uint64]common.Address
	eonKeys    map[uint64][]byte
	chainID    *big.Int
	blocks     map[string]blockData
	txs        map[string]*txtypes.Transaction
	batchIndex uint64
	signer     txtypes.Signer
}

func New(chainID *big.Int, port int16) *SequencerProcessor {
	sequencer := &SequencerProcessor{
		port:       port,
		nonces:     map[string]map[string]uint64{},
		collators:  map[uint64]common.Address{},
		eonKeys:    map[uint64][]byte{},
		chainID:    chainID,
		blocks:     map[string]blockData{},
		txs:        map[string]*txtypes.Transaction{},
		signer:     txtypes.NewLondonSigner(chainID),
		batchIndex: 0,
	}
	// TODO Expose as parameters
	sequencer.setBlock(big.NewInt(22000), 2200000, "latest")
	return sequencer
}

func (proc *SequencerProcessor) setBlock(baseFee *big.Int, gasLimit uint64, block string) {
	b, exists := proc.blocks[block]
	if !exists {
		b = blockData{baseFee: baseFee, gasLimit: gasLimit}
		proc.blocks[block] = b
		return
	}
	b.baseFee = baseFee
	b.gasLimit = gasLimit
}

func (proc *SequencerProcessor) setNonce(a common.Address, nonce uint64, block string) {
	nc, exists := proc.nonces[block]
	if !exists {
		nc = make(map[string]uint64, 0)
		proc.nonces[block] = nc
	}
	nc[a.Hex()] = nonce
}

func (proc *SequencerProcessor) getNonce(a common.Address, block string) uint64 {
	nonce := uint64(0)
	nc, exists := proc.nonces[block]
	if !exists {
		nc = make(map[string]uint64, 0)
		proc.nonces[block] = nc
	}
	nonce, exists = nc[a.Hex()]
	if !exists {
		proc.setNonce(a, nonce, block)
	}
	return nonce
}

func (proc *SequencerProcessor) processEncryptedTx(txBytes []byte) error {
	var tx txtypes.Transaction
	err := tx.UnmarshalBinary(txBytes)
	if err != nil {
		return errors.Wrap(err, "can't unmarshal incoming bytes to transaction")
	}
	if tx.Type() != txtypes.ShutterTxType {
		return errors.New("no shutter tx type")
	}

	sender, err := proc.signer.Sender(&tx)
	if err != nil {
		return errors.New("sender not recoverable")
	}
	nonce := proc.getNonce(sender, "latest")
	if tx.Nonce() != nonce+1 {
		log.Info().Msg("nonce mismatch")
		return nil
	}
	proc.setNonce(sender, nonce+1, "latest")
	return nil
}
