package blockchain

import "strconv"

type Hash [32]byte

func UintToBytes(n uint64) []byte {
	return []byte(strconv.FormatUint(n, 10))
}
