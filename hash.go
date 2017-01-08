package bin

import (
	"crypto/sha256"
	"io"
)

func Hash32(values ...interface{}) uint32 {
	h256 := Hash256(values...)
	return BytesToUint32(h256[:4])
}

func Hash64(values ...interface{}) uint64 {
	h256 := Hash256(values...)
	return BytesToUint64(h256[:8])
}

func Hash128(values ...interface{}) []byte {
	h256 := Hash256(values...)
	return h256[:16]
}

func Hash160(values ...interface{}) []byte {
	h256 := Hash256(values...)
	return h256[:20]
}

func Hash256(values ...interface{}) []byte {
	hash := sha256.New()

	for _, val := range values {
		switch v := val.(type) {

		case []byte:
			hash.Write(v)

		case string:
			hash.Write([]byte(v))

		case int:
			hash.Write(Uint64ToBytes(uint64(v)))

		case uint:
			hash.Write(Uint64ToBytes(uint64(v)))

		case int64:
			hash.Write(Uint64ToBytes(uint64(v)))

		case uint64:
			hash.Write(Uint64ToBytes(v))

		case int32:
			hash.Write(Uint32ToBytes(uint32(v)))

		case uint32:
			hash.Write(Uint32ToBytes(v))

		case error:
			hash.Write([]byte(v.Error()))

		case io.Reader:
			io.Copy(hash, v)

		default:
			panic("SHA256 Not expected value")
		}
	}
	return hash.Sum(nil)
}
