package encoding

import "testing"

func TestStrToBytes(t *testing.T) {
	tests := []struct {
		s string
		b []byte
	}{
		{"", []byte{0, 0, 0, 0, 0, 0, 0, 0}},
		{"a", []byte{0, 0, 0, 0, 0, 0, 0, 1, 97}},
		{"hello", []byte{0, 0, 0, 0, 0, 0, 0, 5, 104, 101, 108, 108, 111}},
	}
	for _, test := range tests {
		b := StrToBytes(test.s)
		if len(b) != 8+len(test.s) {
			t.Errorf("Expected length of %d, but got %d", 8+len(test.s), len(b))
		}
		for i, v := range test.b {
			if b[i] != v {
				t.Errorf("Expected %d, but got %d", v, b[i])
			}
		}
	}
}

func TestStrFromBytes(t *testing.T) {
	tests := []struct {
		s string
		b []byte
	}{
		{"", []byte{0, 0, 0, 0, 0, 0, 0, 0}},
		{"a", []byte{0, 0, 0, 0, 0, 0, 0, 1, 97}},
		{"hello", []byte{0, 0, 0, 0, 0, 0, 0, 5, 104, 101, 108, 108, 111}},
	}
	for _, test := range tests {
		s, l := StrFromBytes(test.b)
		if s != test.s {
			t.Errorf("Expected %s, but got %s", test.s, s)
		}
		if l != uint64(len(test.s)+8) {
			t.Errorf("Expected %d, but got %d", l, uint64(len(test.s)+8))
		}
	}
}
