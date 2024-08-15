package reading

import (
	"github.com/oklog/ulid/v2"
	"livenstore.evrard.online/domain"
)

func ReadEventBytes(b []byte) (domain.Event, []byte) {
	var e domain.Event
	eb, b := ReadBytes(b)
	id, eb := ReadStringBytes(eb)
	ulid, err := ulid.Parse(id)
	if err != nil {
		panic(err)
	}
	e.ID = ulid
	e.Type, eb = ReadStringBytes(eb)
	e.Timestamp, eb = ReadUInt64Bytes(eb)
	e.Data, eb = ReadBytes(eb)
	return e, b
}
