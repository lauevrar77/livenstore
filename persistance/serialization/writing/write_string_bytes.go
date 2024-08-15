package writing

import "livenstore.evrard.online/persistance/serialization/encoding"

func WriteStringBytes(s string, b []byte) []byte {
	return append(b, encoding.StrToBytes(s)...)
}
