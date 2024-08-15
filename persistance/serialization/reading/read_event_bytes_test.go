package reading

import (
	"testing"

	"github.com/oklog/ulid/v2"
	"livenstore.evrard.online/domain"
	"livenstore.evrard.online/persistance/serialization/writing"
)

func TestReadEventBytes(t *testing.T) {
	ulid := ulid.Make()
	tests := []struct {
		input []byte
		event domain.Event
		rest  []byte
	}{
		{writing.WriteEventBytes(domain.Event{ID: ulid, Type: "bar", Timestamp: 42, Data: []byte{1, 2, 3}}, []byte{}), domain.Event{ID: ulid, Type: "bar", Timestamp: 42, Data: []byte{1, 2, 3}}, []byte{}},
		{append(writing.WriteEventBytes(domain.Event{ID: ulid, Type: "bar", Timestamp: 42, Data: []byte{1, 2, 3}}, []byte{}), []byte{1, 2, 3}...), domain.Event{ID: ulid, Type: "bar", Timestamp: 42, Data: []byte{1, 2, 3}}, []byte{1, 2, 3}},
	}

	for _, test := range tests {
		event, rest := ReadEventBytes(test.input)
		if event.ID != test.event.ID {
			t.Errorf("Expected %s, but got %s", test.event.ID, event.ID)
		}
		if event.Type != test.event.Type {
			t.Errorf("Expected %s, but got %s", test.event.Type, event.Type)
		}
		if event.Timestamp != test.event.Timestamp {
			t.Errorf("Expected %d, but got %d", test.event.Timestamp, event.Timestamp)
		}
		if len(event.Data) != len(test.event.Data) {
			t.Errorf("Expected length of %d, but got %d", len(test.event.Data), len(event.Data))
		}
		for i, v := range test.event.Data {
			if event.Data[i] != v {
				t.Errorf("Expected %d, but got %d", v, event.Data[i])
			}
		}
		if len(rest) != len(test.rest) {
			t.Errorf("Expected length of %d, but got %d", len(test.rest), len(rest))
		}
		for i, v := range test.rest {
			if rest[i] != v {
				t.Errorf("Expected %d, but got %d", v, rest[i])
			}
		}
	}
}
