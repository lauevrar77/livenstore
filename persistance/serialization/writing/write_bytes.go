package writing

func WriteBytes(v []byte, b []byte) []byte {

	b = WriteUInt64Bytes(uint64(len(v)), b)

	return append(b, v...)
}
