package replacer_test

import (
	"bytes"
	"testing"

	"github.com/nobishino/replacer"
)

var tests = []struct {
	old, new, in, out string
}{
	{"foo", "bar", "foo", "bar"},
	{"foo", "bar", "ffoo", "fbar"},
	{"foo", "bar", "ffoofooo", "fbarbaro"},
	{"f", "b", "ffoo", "bboo"},
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

// copied from https://github.com/golang/go/blob/72735094660a475a69050b7368c56b25346f5406/src/bytes/bytes_test.go#L1601
type ReplaceTest struct {
	in       string
	old, new string
	n        int
	out      string
}

var ReplaceTests = []ReplaceTest{
	// {"hello", "l", "L", 0, "hello"},
	// {"hello", "l", "L", -1, "heLLo"},
	// {"hello", "x", "X", -1, "hello"},
	// {"", "x", "X", -1, ""},
	// {"radar", "r", "<r>", -1, "<r>ada<r>"},
	// {"", "", "<>", -1, "<>"},
	// {"banana", "a", "<>", -1, "b<>n<>n<>"},
	// {"banana", "a", "<>", 1, "b<>nana"},
	// {"banana", "a", "<>", 1000, "b<>n<>n<>"},
	// {"banana", "an", "<>", -1, "b<><>a"},
	// {"banana", "ana", "<>", -1, "b<>na"},
	// {"banana", "", "<>", -1, "<>b<>a<>n<>a<>n<>a<>"},
	// {"banana", "", "<>", 10, "<>b<>a<>n<>a<>n<>a<>"},
	// {"banana", "", "<>", 6, "<>b<>a<>n<>a<>n<>a"},
	// {"banana", "", "<>", 5, "<>b<>a<>n<>a<>na"},
	// {"banana", "", "<>", 1, "<>banana"},
	// {"banana", "a", "a", -1, "banana"},
	// {"banana", "a", "a", 1, "banana"},
	// {"☺☻☹", "", "<>", -1, "<>☺<>☻<>☹<>"},
}

func TestReplace(t *testing.T) {
	for _, tt := range ReplaceTests {
		in := append([]byte(tt.in), "<spare>"...)
		in = in[:len(tt.in)]
		buf := new(bytes.Buffer)
		w := replacer.NewWriter(buf, []byte(tt.old), []byte(tt.new))
		if _, err := w.Write(in); err != nil {
			t.Fatalf("unexpected error on Write: %v", err)
		}
		if err := w.Flush(); err != nil {
			t.Fatalf("unexpected error on Flush: %v", err)
		}
		if s := buf.String(); s != tt.out {
			t.Errorf("Replace(%q, %q, %q, %d) = %q, want %q", tt.in, tt.old, tt.new, tt.n, s, tt.out)
		}
	}
}
