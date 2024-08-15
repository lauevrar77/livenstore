package writing

import "livenstore.evrard.online/persistance/serialization/encoding"

func WriteUInt64Bytes(n uint64, b []byte) []byte {
	return append(b, encoding.UIntToBytes(n)...)
}
