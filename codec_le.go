//go:build 386 || amd64 || arm || arm64 || wasm || riscv64 || loong64

package vector

import "unsafe"

func encodeFloatsLE(dst []byte, src []float32) {
	if len(src) == 0 {
		return
	}
	srcBytes := unsafe.Slice((*byte)(unsafe.Pointer(&src[0])), len(src)*4)
	copy(dst, srcBytes)
}

func decodeFloatsLE(dst []float32, src []byte) {
	if len(dst) == 0 {
		return
	}
	srcFloats := unsafe.Slice((*float32)(unsafe.Pointer(&src[0])), len(dst))
	copy(dst, srcFloats)
}

func floatsToBytesLE(src []float32) []byte {
	if len(src) == 0 {
		return nil
	}
	return unsafe.Slice((*byte)(unsafe.Pointer(&src[0])), len(src)*4)
}
