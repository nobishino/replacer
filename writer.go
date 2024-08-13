package replacer

import (
	"fmt"
	"io"
)

type Writer struct {
	underlying io.Writer // the Writer to which replaced byte stream is written
	buf        []byte    // internal ring buffer for matching
	filled     int       // number of bytes filled in buf.
	old        []byte    // byte sequence to be replaced
	new        []byte    // byte sequence to replace
}

func (w *Writer) Write(p []byte) (n int, err error) {
	for _, b := range p {
		if err := w.write(b); err != nil {
			return n, err
		}
		n++
	}
	return n, nil
}

// 次に書き込みをするべきbufのindex = 書き込みが完了したidxの次のindex
func (w *Writer) next() int {
	return w.filled % len(w.buf)
}

func (w *Writer) write(b byte) error {
	w.buf[w.next()] = b
	w.filled++
	if w.filled < len(w.buf) {
		return nil
	}
	if w.match() {
		w.filled = 0
		if _, err := w.underlying.Write(w.new); err != nil {
			return err
		}
	} else {
		if _, err := w.underlying.Write(w.buf[:1]); err != nil {
			return err
		}
	}
	return nil
}

func (w *Writer) Flush() error {
	if w.filled > 0 {
		if _, err := w.underlying.Write(w.buf[:w.filled]); err != nil {
			return err
		}
	}
	w.filled = 0
	return nil
}

// 最後に書き込まれたbyteを末尾とする部分列がoldと一致するかどうかを返す
// w.filled < len(w.buf)の場合は呼び出してはいけない(panicする)
func (w *Writer) match() bool {
	if w.filled < len(w.buf) {
		panic(fmt.Sprintf("match should not be called when filled < len(buf). but filled = %d, len(buf) = %d", w.filled, len(w.buf)))
	}
	for i := 0; i < len(w.buf); i++ {
		ii := (w.head() + i) % len(w.buf)
		if w.buf[ii] != w.old[i] {
			return false
		}
	}
	return true
}

// oldとのマッチングを開始するべきindexを返す
func (w *Writer) head() int {
	return w.next()
}

func NewWriter(w io.Writer, old, new []byte) *Writer {
	buf := make([]byte, len(old))
	return &Writer{
		underlying: w,
		buf:        buf,
		old:        []byte(old),
		new:        []byte(new),
	}
}

var _ io.Writer = &Writer{}
