package vector

import "testing"

func TestDot_KnownValues(t *testing.T) {
	// Identical unit vectors -> dot product 1.0
	a := []float32{1.0, 0.0, 0.0}
	b := []float32{1.0, 0.0, 0.0}
	if got := Dot(a, b); got != 1.0 {
		t.Errorf("Dot(identical) = %v, want 1.0", got)
	}

	// Orthogonal vectors -> dot product 0.0
	c := []float32{0.0, 1.0, 0.0}
	if got := Dot(a, c); got != 0.0 {
		t.Errorf("Dot(orthogonal) = %v, want 0.0", got)
	}

	// Hand calculated dot product: 1*2 + 2*3 + 3*4 = 2 + 6 + 12 = 20
	v1 := []float32{1, 2, 3}
	v2 := []float32{2, 3, 4}
	if got := Dot(v1, v2); got != 20.0 {
		t.Errorf("Dot(123, 234) = %v, want 20.0", got)
	}
}

func TestDot_EmptyAndMismatched(t *testing.T) {
	a := []float32{1, 2, 3}
	b := []float32{1, 2}
	if got := Dot(a, b); got != 0 {
		t.Errorf("Dot(mismatched) = %v, want 0", got)
	}

	var empty []float32
	if got := Dot(empty, empty); got != 0 {
		t.Errorf("Dot(empty) = %v, want 0", got)
	}
}

func TestNormalize_UnitLength(t *testing.T) {
	v := []float32{3, 4}
	Normalize(v)
	norm := Norm(v)
	if abs(norm-1.0) > 1e-6 {
		t.Errorf("Norm after Normalize = %v, want ~1.0", norm)
	}
	if abs(v[0]-0.6) > 1e-6 || abs(v[1]-0.8) > 1e-6 {
		t.Errorf("Normalized v = %v, want [0.6, 0.8]", v)
	}
}

func TestNormalize_ZeroVector(t *testing.T) {
	v := []float32{0, 0, 0}
	Normalize(v)
	if v[0] != 0 || v[1] != 0 || v[2] != 0 {
		t.Errorf("Normalize(zero) = %v, want [0, 0, 0]", v)
	}
}

func TestCosine_MatchesDotOfNormalised(t *testing.T) {
	u := []float32{1, 2, 3, 4}
	v := []float32{2, 0, 1, 5}

	cosVal := Cosine(u, v)

	uNorm := append([]float32(nil), u...)
	vNorm := append([]float32(nil), v...)
	Normalize(uNorm)
	Normalize(vNorm)

	dotVal := Dot(uNorm, vNorm)

	if abs(cosVal-dotVal) > 1e-6 {
		t.Errorf("Cosine(%v) = %v, Dot of normalized = %v, difference exceeds 1e-6", cosVal, dotVal, abs(cosVal-dotVal))
	}
}

func BenchmarkDot_384(b *testing.B) {
	v1 := make([]float32, 384)
	v2 := make([]float32, 384)
	for i := 0; i < 384; i++ {
		v1[i] = float32(i) * 0.01
		v2[i] = float32(384-i) * 0.01
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Dot(v1, v2)
	}
}

func abs(x float32) float32 {
	if x < 0 {
		return -x
	}
	return x
}
