package blockchain

import (
	"encoding/hex"
	"strconv"
)

type Hash [32]byte

func (h Hash) String() string {
	return hex.EncodeToString(h[:])
}

func UintToBytes(n uint64) []byte {
	return []byte(strconv.FormatUint(n, 10))
}
