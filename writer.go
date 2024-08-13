package replacer

import "io"

type Writer struct {
	underlying io.Writer
}

func (w *Writer) Write(p []byte) (n int, err error) {
	return w.underlying.Write([]byte("bar"))
}

func NewWriter(w io.Writer, old, new string) *Writer {
	return &Writer{underlying: w}
}
