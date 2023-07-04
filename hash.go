package bin

import "crypto/sha256"

func Hash32(values ...any) uint32 {
	h256 := Hash256(values...)
	return BytesToUint32(h256[:4])
}

func Hash64(values ...any) uint64 {
	h256 := Hash256(values...)
	return BytesToUint64(h256)
}

func Hash128(values ...any) []byte {
	h256 := Hash256(values...)
	return h256[:16]
}

func Hash160(values ...any) []byte {
	h256 := Hash256(values...)
	return h256[:20]
}

func Hash256(values ...any) []byte {
	hash := sha256.New()
	w := NewWriter(hash)
	for _, val := range values {
		w.WriteVar(val)
	}
	return hash.Sum(nil)
}

// FastHash64 is fast non-cryptographic hash function
func FastHash64(values ...any) uint64 {
	data := Encode(values...)
	h := uint64(14695981039346656037)
	for _, c := range data {
		h = (h * 1099511628211) ^ (uint64(c) * 1073676287)
	}
	return h
}
