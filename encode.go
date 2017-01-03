package bin

import "reflect"

type BinEncoder interface {
	BinEncode(*Writer)
}

type BinDecoder interface {
	BinDecode(*Reader)
}

func Encode(obj BinEncoder) ([]byte, error) {
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

func Decode(data []byte, obj BinDecoder) error {
	r := NewBuffer(data)
	if len(data) > 0 {
		//if reflect.ValueOf(obj).IsNil() {
		//	obj = new instance of interface
		//}
		obj.BinDecode(&r.Reader)
	}
	return r.Reader.err
}
