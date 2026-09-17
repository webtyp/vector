//go:build !(386 || amd64 || arm || arm64 || wasm || riscv64 || loong64)

package vector

import "math"

func encodeFloatsLE(dst []byte, src []float32) {
	for i, f := range src {
		bits := math.Float32bits(f)
		dst[i*4+0] = byte(bits)
		dst[i*4+1] = byte(bits >> 8)
		dst[i*4+2] = byte(bits >> 16)
		dst[i*4+3] = byte(bits >> 24)
	}
}

func decodeFloatsLE(dst []float32, src []byte) {
	for i := range dst {
		bits := uint32(src[i*4+0]) | uint32(src[i*4+1])<<8 | uint32(src[i*4+2])<<16 | uint32(src[i*4+3])<<24
		dst[i] = math.Float32frombits(bits)
	}
}

func floatsToBytesLE(src []float32) []byte {
	buf := make([]byte, len(src)*4)
	encodeFloatsLE(buf, src)
	return buf
}
