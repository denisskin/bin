package bin

import "reflect"

type Encoder interface {
	BinEncode(*Writer)
}

type Decoder interface {
	BinDecode(*Reader)
}

func Encode(vv ...interface{}) []byte {
	w := NewBuffer(nil)
	for _, v := range vv {
		w.WriteVar(v)
	}
	return w.Bytes()
}

func Decode(data []byte, vv ...interface{}) error {
	r := NewBuffer(data)
	for _, v := range vv {
		r.ReadVar(v)
	}
	return r.Error()
}

func EncodeObject(obj Encoder) ([]byte, error) {
	w := NewBuffer(nil)
	if obj != nil {
		v := reflect.ValueOf(obj)
		switch v.Kind() {
		case reflect.Chan, reflect.Slice, reflect.Interface, reflect.Ptr, reflect.Map, reflect.Func:
			if !v.IsNil() {
				obj.BinEncode(&w.Writer)
			}
		default:
			obj.BinEncode(&w.Writer)
		}
	}
	return w.Bytes(), w.Writer.err
}

func DecodeObject(data []byte, obj Decoder) error {
	r := NewBuffer(data)
	if len(data) > 0 {
		//if reflect.ValueOf(obj).IsNil() {
		//	obj = new instance of interface
		//}
		obj.BinDecode(&r.Reader)
	}
	return r.Reader.err
}
