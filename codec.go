package vector

import (
	"errors"
	"fmt"
)

var (
	ErrBufferTooShort = errors.New("buffer too short")
	ErrInvalidCount   = errors.New("invalid count or slice size")
)

// ByteLen returns the encoded size of a dim-element vector: dim*4.
func ByteLen(dim int) int {
	if dim < 0 {
		return 0
	}
	return dim * 4
}

// Encode writes v into dst as little-endian float32. dst must be at least
// ByteLen(len(v)). Returns the number of bytes written.
func Encode(dst []byte, v []float32) int {
	needed := ByteLen(len(v))
	if len(dst) < needed {
		return 0
	}
	encodeFloatsLE(dst[:needed], v)
	return needed
}

// Decode reads little-endian float32 from src into dst.
// Returns the number of float32 elements read.
func Decode(dst []float32, src []byte) (int, error) {
	n := len(dst)
	needed := ByteLen(n)
	if len(src) < needed {
		return 0, fmt.Errorf("%w: src length %d < needed %d", ErrBufferTooShort, len(src), needed)
	}
	decodeFloatsLE(dst, src[:needed])
	return n, nil
}

// Bytes returns a zero-copy byte view of the arena's backing store, for handing
// straight to a blob column. Valid until the next Grow.
func (a *Arena) Bytes() []byte {
	return floatsToBytesLE(a.data)
}

// FromBytes fills the arena from an encoded block, in one copy.
func (a *Arena) FromBytes(src []byte, count int) error {
	if a.dim <= 0 {
		return ErrInvalidDimension
	}
	if count < 0 {
		return ErrInvalidCount
	}
	needed := ByteLen(count * a.dim)
	if len(src) < needed {
		return fmt.Errorf("%w: src length %d < needed %d", ErrBufferTooShort, len(src), needed)
	}

	totalFloats := count * a.dim
	if cap(a.data) < totalFloats {
		a.data = make([]float32, totalFloats)
	} else {
		a.data = a.data[:totalFloats]
	}
	a.n = count

	decodeFloatsLE(a.data, src[:needed])
	return nil
}
