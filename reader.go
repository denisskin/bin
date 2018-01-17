package bin

import (
	"encoding/gob"
	"errors"
	"io"
	"math"
	"reflect"
	"time"
)

type Reader struct {
	rd         io.Reader
	err        error
	CntRead    uint64
	MaxCntRead uint64
}

var (
	errBinaryDataWasCorrupted = errors.New("Binary data was corrupted")
	errExceededAllowableLimit = errors.New("Reader.Read-Error: Exceeded allowable limit")
)

func NewReader(rd io.Reader) *Reader {
	return &Reader{rd: rd}
}

func (r *Reader) Error() error {
	return r.err
}

func (r *Reader) ClearError() {
	r.err = nil
}

func (r *Reader) SetError(err error) {
	if err != nil {
		r.err = err
	}
}

func (r *Reader) Close() error {
	if c, ok := r.rd.(io.Closer); ok {
		return c.Close()
	}
	return nil
}

func (r *Reader) SetReadLimit(sz uint64) {
	if sz == 0 {
		r.MaxCntRead = 0
	} else {
		r.MaxCntRead = r.CntRead + sz
	}
}

func (r *Reader) Read(buf []byte) (n int, err error) {
	if r.err != nil {
		return 0, r.err
	}
	defer func() {
		if e, _ := recover().(error); e != nil {
			err = e
		}
		if r.err == nil {
			r.err = err
		}
	}()
	if r.MaxCntRead > 0 && uint64(len(buf))+r.CntRead > r.MaxCntRead {
		err = errExceededAllowableLimit
		return
	}
	n, err = io.ReadFull(r.rd, buf)
	r.CntRead += uint64(n)
	return
}

func (r *Reader) read(length int) ([]byte, error) {
	buf := make([]byte, length)
	_, err := r.Read(buf)
	return buf, err
}

//----------- fixed types --------------
func (r *Reader) ReadUint8() (uint8, error) {
	bb, err := r.read(1)
	if len(bb) == 1 {
		return uint8(bb[0]), err
	}
	return 0, err
}

func (r *Reader) ReadUint16() (uint16, error) {
	b, err := r.read(2)
	return BytesToUint16(b), err
}

func (r *Reader) ReadUint32() (uint32, error) {
	b, err := r.read(4)
	return BytesToUint32(b), err
}

func (r *Reader) ReadUint64() (uint64, error) {
	b, err := r.read(8)
	return BytesToUint64(b), err
}

func (r *Reader) ReadFloat32() (float32, error) {
	b, err := r.read(4)
	return math.Float32frombits(BytesToUint32(b)), err
}

func (r *Reader) ReadFloat64() (float64, error) {
	b, err := r.read(8)
	return math.Float64frombits(BytesToUint64(b)), err
}

func (r *Reader) ReadTime() (time.Time, error) {
	v, err := r.ReadUint64()
	return time.Unix(0, int64(v)), err
}

func (r *Reader) ReadTime32() (time.Time, error) {
	v, err := r.ReadUint32()
	return time.Unix(int64(v), 0), err
}

func (r *Reader) ReadByte() (b byte, err error) {
	bb, err := r.read(1)
	if len(bb) > 0 {
		b = bb[0]
	}
	return
}

func (r *Reader) ReadBool() (bool, error) {
	b, err := r.ReadByte()
	return b != 0, err
}

//----------- var types ----------------
func (r *Reader) readVarInt() (i int64) {
	b0, err := r.ReadUint8()
	if err != nil {
		return
	}
	if b0&0x80 == 0 {
		return int64(b0)
	}
	n := int(b0 & 0x0f)
	if n > 0 {
		if n > 8 {
			r.SetError(errBinaryDataWasCorrupted)
			return
		}
		bb, err := r.read(n)
		if err != nil {
			return
		}
		for _, c := range bb {
			i <<= 8
			i |= int64(c)
		}
		if b0&0x40 != 0 {
			i = -i
		}
	}
	return
}

func (r *Reader) ReadVarInt() (int, error) {
	v := r.readVarInt()
	return int(v), r.err
}

func (r *Reader) ReadVarUint() (uint64, error) {
	v := r.readVarInt()
	return uint64(v), r.err
}

func (r *Reader) ReadSliceBytes() ([][]byte, error) {
	n, err := r.ReadVarInt()
	if err != nil {
		return nil, err
	}
	res := make([][]byte, n)
	for i := 0; i < n; i++ {
		if res[i], err = r.ReadBytes(); err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (r *Reader) ReadBytes() ([]byte, error) {
	if n, err := r.ReadVarInt(); err != nil {
		return nil, err
	} else if n > 0 {
		return r.read(n)
	}
	return []byte{}, nil
}

func (r *Reader) ReadString() (string, error) {
	v, err := r.ReadBytes()
	return string(v), err
}

func (r *Reader) ReadSliceString() ([]string, error) {
	n, err := r.ReadVarInt()
	if err != nil {
		return nil, err
	}
	res := make([]string, n)
	for i := 0; i < n; i++ {
		if res[i], err = r.ReadString(); err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (r *Reader) ReadError() (error, error) {
	if s, err := r.ReadString(); err != nil {
		return nil, err
	} else {
		return errors.New(s), nil
	}
}

func (r *Reader) ReadVar(val interface{}) error {
	switch v := val.(type) {
	case *int:
		*v = int(r.readVarInt())
	case *int8:
		*v = int8(r.readVarInt())
	case *int16:
		*v = int16(r.readVarInt())
	case *int32:
		*v = int32(r.readVarInt())
	case *int64:
		*v = int64(r.readVarInt())

	case *uint:
		*v = uint(r.readVarInt())
	case *uint8:
		*v = uint8(r.readVarInt())
	case *uint16:
		*v = uint16(r.readVarInt())
	case *uint32:
		*v = uint32(r.readVarInt())
	case *uint64:
		*v = uint64(r.readVarInt())

	case *float32:
		*v, _ = r.ReadFloat32()
	case *float64:
		*v, _ = r.ReadFloat64()
	case *time.Time:
		*v, _ = r.ReadTime()
	case *bool:
		*v, _ = r.ReadBool()

	case *string:
		*v, _ = r.ReadString()
	case *[]string:
		*v, _ = r.ReadSliceString()
	case *[]byte:
		*v, _ = r.ReadBytes()
	case *[][]byte:
		*v, _ = r.ReadSliceBytes()

	case BinDecoder:
		bb, err := r.ReadBytes()
		if err != nil {
			return err
		}
		if err := DecodeObject(bb, v); err != nil {
			r.SetError(err)
		}

	case *error:
		*v, _ = r.ReadError()

	default:

		// read object in case:  var obj*Object; r.Read(&obj)
		if pp := reflect.ValueOf(val); pp.Kind() == reflect.Ptr && !pp.IsNil() {
			if p := pp.Elem(); p.Kind() == reflect.Ptr { //  && p.IsNil()
				objPtr := reflect.New(reflect.TypeOf(p.Interface()).Elem())
				if obj, ok := objPtr.Interface().(BinDecoder); ok {
					if bb, _ := r.ReadBytes(); len(bb) > 0 {
						if err := DecodeObject(bb, obj); err != nil {
							r.SetError(err)
						} else {
							p.Set(objPtr)
						}
					}
					break
				}
			}
		}

		//case reflect.Chan, reflect.Slice, reflect.Interface, reflect.Ptr, reflect.Map, reflect.Func:
		//	if !v.IsNil() {
		//		obj.BinEncode(w)
		//	}
		//}

		// other type
		r.err = gob.NewDecoder(r).Decode(v)
	}
	return r.err
}
