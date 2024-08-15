package encoding

func StrToBytes(s string) []byte {
	b := []byte(s)
	lb := UIntToBytes(uint64(len(b)))
	return append(lb, b...)
}

func StrFromBytes(b []byte) (string, uint64) {
	l := UInt64FromBytes(b[:8])
	return string(b[8 : 8+l]), 8 + l
}
