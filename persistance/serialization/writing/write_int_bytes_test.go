package writing

import (
	"testing"

	"livenstore.evrard.online/persistance/serialization/encoding"
)

func TestWriteUInt64Bytes(t *testing.T) {
	tests := []struct {
		before   []byte
		value    uint64
		expected []byte
	}{
		{[]byte{}, 0, encoding.UIntToBytes(0)},
		{[]byte{1, 2, 3}, 4, append([]byte{1, 2, 3}, encoding.UIntToBytes(4)...)},
	}

	for _, test := range tests {
		b := WriteUInt64Bytes(test.value, test.before)
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
