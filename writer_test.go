package replacer_test

import (
	"bytes"
	"testing"

	"github.com/nobishino/replacer"
)

// func TestWriter(t *testing.T) {
// 	buf := new(bytes.Buffer)
// 	w := replacer.NewWriter(buf, "foo", "bar")
// 	w.Write([]byte("foo"))
// 	if buf.String() != "bar" {
// 		t.Errorf("got %q, want %q", buf.String(), "bar")
// 	}
// }

var tests = []struct {
	old, new, in, out string
}{
	{"foo", "bar", "foo", "bar"},
	{"foo", "bar", "ffoo", "fbar"},
	{"foo", "bar", "ffoofooo", "fbarbaro"},
}

func TestWriter(t *testing.T) {
	for _, tt := range tests {
		buf := new(bytes.Buffer)
		w := replacer.NewWriter(buf, []byte(tt.old), []byte(tt.new))

		if _, err := w.Write([]byte(tt.in)); err != nil {
			t.Errorf("unexpected error on Write: %v", err)
		}
		if err := w.Flush(); err != nil {
			t.Fatalf("unexpected error on Flush: %v", err)
		}
		if buf.String() != tt.out {
			t.Errorf("got %q, want %q", buf.String(), tt.out)
		}
	}
}
