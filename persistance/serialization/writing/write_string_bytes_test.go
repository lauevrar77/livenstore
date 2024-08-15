package writing

import (
	"testing"

	"livenstore.evrard.online/persistance/serialization/encoding"
)

func TestWriteStringBytes(t *testing.T) {
	tests := []struct {
		before   []byte
		value    string
		expected []byte
	}{
		{[]byte{}, "foo", encoding.StrToBytes("foo")},
		{[]byte{1, 2, 3}, "bar", append([]byte{1, 2, 3}, encoding.StrToBytes("bar")...)},
	}

	for _, test := range tests {
		b := WriteStringBytes(test.value, test.before)
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
