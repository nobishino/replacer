package replacer

import "io"

type Writer struct {
	underlying io.Writer
	buf        []byte
}

func (w *Writer) Write(p []byte) (n int, err error) {
	return w.underlying.Write([]byte("bar"))
}

func NewWriter(w io.Writer, old, new string) *Writer {
	buf := make([]byte, len(old))
	return &Writer{
		underlying: w,
		buf:        buf,
	}
}
