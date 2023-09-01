// Copyright (C) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package database

import "github.com/ava-labs/avalanchego/ids"

const (
	LatestSeenBlockKey = "latestSeenBlock"
)

// RelayerDatabase is a key-value store for relayer state, with each chainID maintaining its own state
type RelayerDatabase interface {
	Get(chainID ids.ID, key []byte) ([]byte, error)
	Put(chainID ids.ID, key []byte, value []byte) error
}
