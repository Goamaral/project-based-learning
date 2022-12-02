package blockchain

import "strconv"

func IntToBytes(n int64) []byte {
	return []byte(strconv.FormatInt(n, 10))
}
