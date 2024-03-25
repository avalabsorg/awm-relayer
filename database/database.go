// Copyright (C) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

//go:generate mockgen -source=$GOFILE -destination=./mocks/mock_database.go -package=mocks

package database

import (
	"strings"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/awm-relayer/config"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
)

const (
	LatestProcessedBlockKey = "latestProcessedBlock"
)

var (
	ErrDataKeyNotFound          = errors.New("data key not found")
	ErrRelayerKeyNotFound       = errors.New("no database for relayer key")
	ErrDatabaseMisconfiguration = errors.New("database misconfiguration")
)

// RelayerDatabase is a key-value store for relayer state, with each relayerKey maintaining its own state
type RelayerDatabase interface {
	Get(relayerKey common.Hash, dataKey []byte) ([]byte, error)
	Put(relayerKey common.Hash, dataKey []byte, value []byte) error
}

// Returns true if an error returned by a RelayerDatabase indicates the requested key was not found
func IsKeyNotFoundError(err error) bool {
	return errors.Is(err, ErrRelayerKeyNotFound) || errors.Is(err, ErrDataKeyNotFound)
}

// RelayerKey is a unique identifier for an application relayer
type RelayerKey struct {
	SourceBlockchainID      ids.ID
	DestinationBlockchainID ids.ID
	OriginSenderAddress     common.Address
	DestinationAddress      common.Address
}

// CalculateRelayerKey calculates the unique identifier for an application relayer
func (k RelayerKey) CalculateRelayerKey() common.Hash {
	return CalculateRelayerKey(
		k.SourceBlockchainID,
		k.DestinationBlockchainID,
		k.OriginSenderAddress,
		k.DestinationAddress,
	)
}

// Standalone utility to calculate a relayer key
func CalculateRelayerKey(
	sourceBlockchainID ids.ID,
	destinationBlockchainID ids.ID,
	originSenderAddress common.Address,
	desinationAddress common.Address,
) common.Hash {
	return crypto.Keccak256Hash(
		[]byte(strings.Join(
			[]string{
				sourceBlockchainID.String(),
				destinationBlockchainID.String(),
				originSenderAddress.String(),
				desinationAddress.String(),
			},
			"-",
		)),
	)
}

// Get all of the possible relayer keys for a given configuration
func GetConfigRelayerKeys(cfg *config.Config) []RelayerKey {
	var keys []RelayerKey
	for _, s := range cfg.SourceBlockchains {
		keys = append(keys, GetSourceConfigRelayerKeys(s)...)
	}
	return keys
}

// Calculate all of the possible relayer keys for a given source blockchain
func GetSourceConfigRelayerKeys(cfg *config.SourceBlockchain) []RelayerKey {
	var keys []RelayerKey
	for _, dst := range cfg.GetSupportedDestinations().List() {
		keys = append(keys, RelayerKey{
			SourceBlockchainID:      cfg.GetBlockchainID(),
			DestinationBlockchainID: dst,
			OriginSenderAddress:     common.Address{}, // TODO: populate with allowed sender/receiver addresses
			DestinationAddress:      common.Address{},
		})
	}
	return keys
}
