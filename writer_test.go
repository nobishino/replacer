package replacer_test

import (
	"bytes"
	"testing"

	"github.com/nobishino/replacer"
)

func TestWriter(t *testing.T) {
	buf := new(bytes.Buffer)
	w := replacer.NewWriter(buf, "foo", "bar")
	w.Write([]byte("foo"))
	if buf.String() != "bar" {
		t.Errorf("got %q, want %q", buf.String(), "bar")
	}
}
