package writing

import (
	"livenstore.evrard.online/persistance/serialization/encoding"
	"testing"
)

func TestWriteBytes(t *testing.T) {
	tests := []struct {
		v        []byte
		b        []byte
		expected []byte
	}{
		{[]byte{}, []byte{0, 1, 2, 3, 4, 5, 6, 7}, append(append([]byte{0, 1, 2, 3, 4, 5, 6, 7}, encoding.UIntToBytes(uint64(0))...), []byte{}...)},
		{[]byte{1, 2, 3}, []byte{0, 1, 2, 3, 4, 5, 6, 7}, append(append([]byte{0, 1, 2, 3, 4, 5, 6, 7}, encoding.UIntToBytes(uint64(3))...), []byte{1, 2, 3}...)},
	}

	for _, test := range tests {
		b := WriteBytes(test.v, test.b)
		if len(b) != len(test.expected) {
			t.Errorf("Expected length of %d, but got %d", len(test.expected), len(b))
		}
		for i, v := range test.expected {
			if b[i] != v {
				t.Errorf("Expected %d, but got %d", v, b[i])
			}
		}
	}
}
