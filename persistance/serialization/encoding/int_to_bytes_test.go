package encoding

import "testing"

func TestUIntToBytes(t *testing.T) {
	tests := []struct {
		n uint64
		b []byte
	}{
		{0, []byte{0, 0, 0, 0, 0, 0, 0, 0}},
		{1, []byte{0, 0, 0, 0, 0, 0, 0, 1}},
		{256, []byte{0, 0, 0, 0, 0, 0, 1, 0}},
		{18446744073709551615, []byte{255, 255, 255, 255, 255, 255, 255, 255}},
	}

	for _, test := range tests {
		b := UIntToBytes(test.n)
		if len(b) != 8 {
			t.Errorf("Expected length of 8, but got %d", len(b))
		}
		for i, v := range test.b {
			if b[i] != v {
				t.Errorf("Expected %d, but got %d", v, b[i])
			}
		}
	}

}

func TestUInt64FromBytes(t *testing.T) {
	tests := []struct {
		n uint64
		b []byte
	}{
		{0, []byte{0, 0, 0, 0, 0, 0, 0, 0}},
		{1, []byte{0, 0, 0, 0, 0, 0, 0, 1}},
		{256, []byte{0, 0, 0, 0, 0, 0, 1, 0}},
		{18446744073709551615, []byte{255, 255, 255, 255, 255, 255, 255, 255}},
	}

	for _, test := range tests {
		n := UInt64FromBytes(test.b)
		if n != test.n {
			t.Errorf("Expected %d, but got %d", test.n, n)
		}
	}
}
