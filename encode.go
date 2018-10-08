package bin

type Encoder interface {
	Encode() []byte
}

type Decoder interface {
	Decode([]byte) error
}

func Encode(vv ...interface{}) []byte {
	w := NewBuffer(nil)
	w.WriteVar(vv...)
	return w.Bytes()
}

func Decode(data []byte, vv ...interface{}) error {
	return NewBuffer(data).ReadVar(vv...)
}

type binWriter interface {
	BinWrite(writer *Writer)
}

type binReader interface {
	BinRead(reader *Reader)
}
