package encoding

import "encoding/binary"

func UIntToBytes(n uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, n)
	return b
}

func UInt64FromBytes(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}
