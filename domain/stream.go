package domain

import "github.com/oklog/ulid/v2"

type Stream struct {
	Name     string
	EventIDs []ulid.ULID
}
