package reading

func ReadBytes(b []byte) ([]byte, []byte) {
	l, b := ReadUInt64Bytes(b)
	return b[:l], b[l:]
}
