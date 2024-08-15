package reading

import (
	"livenstore.evrard.online/persistance/serialization/encoding"
	"testing"
)

func TestReadUInt64Bytes(t *testing.T) {
	tests := []struct {
		b []byte
		n uint64
		r []byte
	}{
		{append(encoding.UIntToBytes(0), []byte{1, 2, 3}...), 0, []byte{1, 2, 3}},
		{append(encoding.UIntToBytes(1), []byte{1, 2, 3}...), 1, []byte{1, 2, 3}},
		{append(encoding.UIntToBytes(256), []byte{1, 2, 3}...), 256, []byte{1, 2, 3}},
		{append(encoding.UIntToBytes(18446744073709551615), []byte{1, 2, 3, 4, 5}...), 18446744073709551615, []byte{1, 2, 3, 4, 5}},
	}

	for _, test := range tests {
		n, r := ReadUInt64Bytes(test.b)
		if n != test.n {
			t.Errorf("Expected %d, but got %d", test.n, n)
		}
		if len(r) != len(test.r) {
			t.Errorf("Expected length of %d, but got %d", len(test.r), len(r))
		}
		for i, v := range test.r {
			if r[i] != v {
				t.Errorf("Expected %d, but got %d", v, r[i])
			}
		}
	}
}
