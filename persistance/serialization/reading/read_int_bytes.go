package reading

import "livenstore.evrard.online/persistance/serialization/encoding"

func ReadUInt64Bytes(b []byte) (uint64, []byte) {
	return encoding.UInt64FromBytes(b[:8]), b[8:]
}
