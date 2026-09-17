//go:build !(386 || amd64 || arm || arm64 || wasm || riscv64 || loong64)

package vector

import (
	"encoding/binary"
	"math"
)

func encodeFloatsLE(dst []byte, src []float32) {
	for i, f := range src {
		binary.LittleEndian.PutUint32(dst[i*4:], math.Float32bits(f))
	}
}

func decodeFloatsLE(dst []float32, src []byte) {
	for i := range dst {
		dst[i] = math.Float32frombits(binary.LittleEndian.Uint32(src[i*4:]))
	}
}

func floatsToBytesLE(src []float32) []byte {
	buf := make([]byte, len(src)*4)
	encodeFloatsLE(buf, src)
	return buf
}
