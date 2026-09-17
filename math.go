package vector

import "math"

// Dot returns the dot product of two vectors.
// For L2-normalised vectors this IS the cosine similarity.
// Returns 0 if lengths are mismatched or vectors are empty.
func Dot(a, b []float32) float32 {
	if len(a) != len(b) || len(a) == 0 {
		return 0
	}
	var sum float32
	n := len(a)
	i := 0
	for ; i <= n-4; i += 4 {
		sum += a[i]*b[i] + a[i+1]*b[i+1] + a[i+2]*b[i+2] + a[i+3]*b[i+3]
	}
	for ; i < n; i++ {
		sum += a[i] * b[i]
	}
	return sum
}

// Norm returns the L2 magnitude (Euclidean norm) of vector v.
func Norm(v []float32) float32 {
	if len(v) == 0 {
		return 0
	}
	var sum float64
	for _, x := range v {
		sum += float64(x) * float64(x)
	}
	return float32(math.Sqrt(sum))
}

// Normalize rescales v in place to have unit length (L2 norm = 1).
// A zero or empty vector is left untouched.
func Normalize(v []float32) {
	norm := Norm(v)
	if norm == 0 || math.IsNaN(float64(norm)) {
		return
	}
	inv := 1.0 / norm
	for i := range v {
		v[i] *= inv
	}
}

// Cosine returns the cosine similarity between two un-normalised vectors.
// For L2-normalised vectors, Dot is preferred.
func Cosine(a, b []float32) float32 {
	if len(a) != len(b) || len(a) == 0 {
		return 0
	}
	normA := Norm(a)
	normB := Norm(b)
	if normA == 0 || normB == 0 {
		return 0
	}
	return Dot(a, b) / (normA * normB)
}
