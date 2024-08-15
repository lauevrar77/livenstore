package reading

import (
	"testing"

	"livenstore.evrard.online/persistance/serialization/encoding"
)

func TestReadStringBytes(t *testing.T) {
	tests := []struct {
		input []byte
		value string
		rest  []byte
	}{
		{encoding.StrToBytes("foo"), "foo", []byte{}},
		{append(encoding.StrToBytes("foo"), 1, 2, 3), "foo", []byte{1, 2, 3}},
	}

	for _, test := range tests {
		value, rest := ReadStringBytes(test.input)
		if value != test.value {
			t.Errorf("Expected %s, but got %s", test.value, value)
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
