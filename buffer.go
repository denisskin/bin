package bin

import "bytes"

type Buffer struct {
	buf *bytes.Buffer
	Reader
	Writer
}

func NewBuffer(p []byte, v ...interface{}) *Buffer {
	b := bytes.NewBuffer(p)
	buf := &Buffer{
		b,
		Reader{rd: b},
		Writer{wr: b},
	}
	if len(v) > 0 {
		buf.WriteVar(v...)
	}
	return buf
}

func (b *Buffer) Error() error {
	if b.Reader.err != nil {
		return b.Reader.err
	}
	return b.Writer.err
}

func (w *Buffer) Bytes() []byte {
	return w.buf.Bytes()
}

func (w *Buffer) Buffer() *bytes.Buffer {
	return w.buf
}
