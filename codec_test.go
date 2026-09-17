package vector

import (
	"bytes"
	"math/rand"
	"testing"
)

func TestCodec_RoundTrip(t *testing.T) {
	rnd := rand.New(rand.NewSource(42))
	floats := make([]float32, 384)
	for i := range floats {
		floats[i] = rnd.Float32()
	}

	buf := make([]byte, ByteLen(len(floats)))
	nEnc := Encode(buf, floats)
	if nEnc != len(buf) {
		t.Fatalf("Encode wrote %d bytes, want %d", nEnc, len(buf))
	}

	decoded := make([]float32, len(floats))
	nDec, err := Decode(decoded, buf)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	if nDec != len(floats) {
		t.Fatalf("Decode read %d floats, want %d", nDec, len(floats))
	}

	for i := range floats {
		if floats[i] != decoded[i] {
			t.Errorf("At index %d: encoded %v, decoded %v", i, floats[i], decoded[i])
		}
	}
}

func TestCodec_LittleEndianLayout(t *testing.T) {
	// Encode([1.0]) produce exactly 00 00 80 3F
	v := []float32{1.0}
	buf := make([]byte, 4)
	Encode(buf, v)

	want := []byte{0x00, 0x00, 0x80, 0x3f}
	if !bytes.Equal(buf, want) {
		t.Errorf("Encode([1.0]) = %x, want %x", buf, want)
	}
}

func TestCodec_DecodeShortInput(t *testing.T) {
	dst := make([]float32, 4)
	shortBuf := make([]byte, 10) // needs 16 bytes for 4 floats

	_, err := Decode(dst, shortBuf)
	if err == nil {
		t.Error("Decode short input should return error, got nil")
	}
}

func TestArena_BytesFromBytesRoundTrip(t *testing.T) {
	dim := 128
	count := 1024
	rnd := rand.New(rand.NewSource(123))

	arena1 := NewArena(dim, count)
	v := make([]float32, dim)
	for i := 0; i < count; i++ {
		for j := range v {
			v[j] = rnd.Float32()
		}
		arena1.Append(v)
	}

	blob := arena1.Bytes()

	arena2 := NewArena(dim, 0)
	err := arena2.FromBytes(blob, count)
	if err != nil {
		t.Fatalf("FromBytes failed: %v", err)
	}

	if arena2.Len() != count {
		t.Fatalf("arena2 Len = %d, want %d", arena2.Len(), count)
	}

	for i := 0; i < count; i++ {
		v1 := arena1.At(i)
		v2 := arena2.At(i)
		for j := 0; j < dim; j++ {
			if v1[j] != v2[j] {
				t.Fatalf("Mismatch at vector %d, dim %d: %v vs %v", i, j, v1[j], v2[j])
			}
		}
	}
}
