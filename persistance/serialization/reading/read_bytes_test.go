package reading

import (
	"livenstore.evrard.online/persistance/serialization/writing"
	"testing"
)

func TestReadBytes(t *testing.T) {
	tests := []struct {
		b     []byte
		value []byte
		rest  []byte
	}{
		{writing.WriteBytes([]byte{1, 2, 3}, []byte{}), []byte{1, 2, 3}, []byte{}},
		{append(writing.WriteBytes([]byte{1, 2, 3}, []byte{}), []byte{4, 5, 6}...), []byte{1, 2, 3}, []byte{4, 5, 6}},
	}

	for _, test := range tests {
		value, rest := ReadBytes(test.b)
		if len(value) != len(test.value) {
			t.Errorf("Expected length of %d, but got %d", len(test.value), len(value))
		}
		for i, v := range test.value {
			if value[i] != v {
				t.Errorf("Expected %d, but got %d", v, value[i])
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
