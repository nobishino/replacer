package replacer

import "io"

type Writer struct {
	underlying io.Writer // the Writer to which replaced byte stream is written
	buf        []byte    // internal buffer for matching
	filled     int       // buf[x] == 0 の時にそれが書き込みされていないのか0が書き込まれたのかを区別するためのindex
	old        []byte
	new        []byte
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
	w.buf[w.filled%len(w.buf)] = b
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

// 最後に書き込まれたbyteを末尾とする部分列がoldと一致するかどうかを返す
func (w *Writer) match() bool {
	if w.filled < len(w.buf) {
		return false
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
