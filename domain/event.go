package domain

import "github.com/oklog/ulid/v2"

type Event struct {
	ID        ulid.ULID
	Type      string
	Timestamp uint64
	Data      []byte
}
