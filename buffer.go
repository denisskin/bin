package bin

import (
	"bytes"
)

type Buffer struct {
	buf *bytes.Buffer
	Reader
	Writer
}

func NewBuffer(p []byte) *Buffer {
	buf := &Buffer{
		buf: bytes.NewBuffer(p),
	}
	buf.Reader = *NewReader(buf.buf)
	buf.Writer = *NewWriter(buf.buf)
	return buf
}

func (b *Buffer) Error() error {
	if b.Reader.err != nil {
		return b.Reader.err
	}
	return b.Reader.err
}

func (w *Buffer) Bytes() []byte {
	return w.buf.Bytes()
}
