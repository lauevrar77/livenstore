package writing

import "livenstore.evrard.online/domain"

func WriteEventBytes(e domain.Event, b []byte) []byte {
	eb := make([]byte, 0)
	eb = WriteStringBytes(e.ID.String(), eb)
	eb = WriteStringBytes(e.Type, eb)
	eb = WriteUInt64Bytes(e.Timestamp, eb)
	eb = WriteBytes(e.Data, eb)
	return WriteBytes(eb, b)
}
