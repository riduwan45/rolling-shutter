package shutterevents

import (
	"crypto/ecdsa"
	"encoding/base64"
	"encoding/hex"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/bls12381"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/pkg/errors"

	"github.com/shutter-network/shutter/shlib/shcrypto"
)

func encodeUint64(val uint64) string {
	return strconv.FormatUint(val, 10)
}

func decodeUint64(val string) (uint64, error) {
	v, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return 0, errors.Wrap(err, "failed to parse event")
	}
	return v, nil
}

// encodeAddresses encodes the given slice of Addresses as comma-separated list of addresses.
func encodeAddresses(addr []common.Address) string {
	var hexstrings []string
	for _, a := range addr {
		hexstrings = append(hexstrings, a.Hex())
	}
	return strings.Join(hexstrings, ",")
}

// decodeAddresses reverses the encodeAddressesForEvent operation, i.e. it parses a list
// of addresses from a comma-separated string.
func decodeAddresses(s string) ([]common.Address, error) {
	var res []common.Address
	if s == "" {
		return res, nil
	}
	for _, a := range strings.Split(s, ",") {
		if !common.IsHexAddress(a) {
			return nil, errors.Errorf("malformed address: %q", s)
		}

		res = append(res, common.HexToAddress(a))
	}
	return res, nil
}

// encodeByteSequence encodes a slice o byte strings as a comma separated string.
func encodeByteSequence(v [][]byte) string {
	var hexstrings []string
	for _, a := range v {
		hexstrings = append(hexstrings, hexutil.Encode(a))
	}
	return strings.Join(hexstrings, ",")
}

// decodeByteSequence parses a list of hex encoded, comma-separated byte slices.
func decodeByteSequence(s string) ([][]byte, error) {
	var res [][]byte
	if s == "" {
		return res, nil
	}
	for _, v := range strings.Split(s, ",") {
		bs, err := hexutil.Decode(v)
		if err != nil {
			return [][]byte{}, err
		}
		res = append(res, bs)
	}
	return res, nil
}

// encodePubkey encodes the PublicKey as a string suitable for putting it into a tendermint
// event, i.e. an utf-8 compatible string.
func encodePubkey(pubkey *ecdsa.PublicKey) string {
	return base64.RawURLEncoding.EncodeToString(ethcrypto.FromECDSAPub(pubkey))
}

// decodePubkey decodes a public key from a tendermint event.
func decodePubkey(val string) (*ecdsa.PublicKey, error) {
	data, err := base64.RawURLEncoding.DecodeString(val)
	if err != nil {
		return nil, err
	}
	return ethcrypto.UnmarshalPubkey(data)
}

func encodeGammas(gammas *shcrypto.Gammas) string {
	g2 := bls12381.NewG2()
	var encoded []string
	if gammas != nil {
		for _, g := range *gammas {
			encoded = append(encoded, hex.EncodeToString(g2.ToBytes(g)))
		}
	}
	return strings.Join(encoded, ",")
}

func decodeGammas(eventValue string) (shcrypto.Gammas, error) {
	g2 := bls12381.NewG2()
	parts := strings.Split(eventValue, ",")
	var res shcrypto.Gammas
	for _, p := range parts {
		marshaledG2, err := hex.DecodeString(p)
		if err != nil {
			return shcrypto.Gammas{}, err
		}
		g, err := g2.FromBytes(marshaledG2)
		if err != nil {
			return shcrypto.Gammas{}, err
		}
		if !g2.IsOnCurve(g) {
			return shcrypto.Gammas{}, errors.Errorf("invalid gamma value %x", p)
		}
		res = append(res, g)
	}
	return res, nil
}

func encodeAddress(a common.Address) string {
	return a.Hex()
}

func decodeAddress(s string) (common.Address, error) {
	a := common.HexToAddress(s)
	if a.Hex() != s {
		return common.Address{}, errors.Errorf("invalid address %s", s)
	}
	return a, nil
}

func encodeECIESPublicKey(key *ecies.PublicKey) string {
	return encodePubkey(key.ExportECDSA())
}

func decodeECIESPublicKey(val string) (*ecies.PublicKey, error) {
	publicKeyECDSA, err := decodePubkey(val)
	if err != nil {
		return nil, err
	}
	return ecies.ImportECDSAPublic(publicKeyECDSA), nil
}
