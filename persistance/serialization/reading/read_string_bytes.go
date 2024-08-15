package reading

import "livenstore.evrard.online/persistance/serialization/encoding"

func ReadStringBytes(b []byte) (string, []byte) {
	s, l := encoding.StrFromBytes(b)
	return s, b[l:]
}
