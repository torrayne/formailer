package formailer

import (
	"strings"
	"testing"
)

func TestIsSlice(t *testing.T) {
	tests := map[string]bool{
		"":          false,
		"not_slice": false,
		"a slice":   true,
	}

	for str, expected := range tests {
		var r bool
		split := strings.Split(str, " ")
		if len(split) > 1 {
			r = isSlice(split)
		} else {
			r = isSlice(split[0])
		}

		if r != expected {
			t.Errorf("Unexpected result from isSlice. On: %s\nExpected: %t; Got: %t", str, expected, r)
		}
	}
}
