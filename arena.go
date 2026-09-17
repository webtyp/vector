package vector

import "webtyp.com/fmt"

// Sentinel errors for the simple, detail-free cases — matching the convention
// used across storage.ErrNoRows / orm.ErrNotFound / rbac.ErrNotFound.
var (
	ErrDimensionMismatch = fmt.Err("vector: dimension mismatch")
	ErrIndexOutOfBounds  = fmt.Err("vector: index out of bounds")
	ErrInvalidDimension  = fmt.Err("vector: dimension must be greater than zero")
)

// Arena is a contiguous block of N × Dim float32 values. Vector i occupies
// data[i*dim : (i+1)*dim].
//
// One allocation holds every vector in the corpus. This is not a micro-optimisation:
// TinyGo's conservative GC scans N small slices far more expensively than one large
// one, and a contiguous layout is what makes the dot product walk memory linearly.
type Arena struct {
	data []float32
	dim  int
	n    int
}

// NewArena initializes an Arena with a given dimension and capacity.
func NewArena(dim, capacity int) *Arena {
	if dim <= 0 {
		dim = 0
	}
	if capacity < 0 {
		capacity = 0
	}
	return &Arena{
		data: make([]float32, 0, dim*capacity),
		dim:  dim,
		n:    0,
	}
}

// Dim returns the vector dimension of the arena.
func (a *Arena) Dim() int {
	return a.dim
}

// Len returns the number of vectors stored in the arena.
func (a *Arena) Len() int {
	return a.n
}

// Cap returns the total vector capacity of the arena before reallocating.
func (a *Arena) Cap() int {
	if a.dim == 0 {
		return 0
	}
	return cap(a.data) / a.dim
}

// Append copies v into the next slot, normalising it in place, and returns its index.
// v is not retained.
func (a *Arena) Append(v []float32) (int, error) {
	if a.dim <= 0 {
		return -1, ErrInvalidDimension
	}
	if len(v) != a.dim {
		return -1, fmt.Err("vector: dimension mismatch, expected", a.dim, "got", len(v))
	}

	idx := a.n
	start := idx * a.dim
	end := start + a.dim

	if cap(a.data) < end {
		a.Grow(a.n + 1)
	}

	a.data = a.data[:end]
	copy(a.data[start:end], v)
	Normalize(a.data[start:end])
	a.n++

	return idx, nil
}

// At returns a view into the arena — NOT a copy. Writing through it mutates the arena.
func (a *Arena) At(i int) []float32 {
	if i < 0 || i >= a.n {
		return nil
	}
	start := i * a.dim
	return a.data[start : start+a.dim]
}

// Set overwrites slot i, normalising.
func (a *Arena) Set(i int, v []float32) error {
	if i < 0 || i >= a.n {
		return ErrIndexOutOfBounds
	}
	if len(v) != a.dim {
		return fmt.Err("vector: dimension mismatch, expected", a.dim, "got", len(v))
	}

	start := i * a.dim
	copy(a.data[start:start+a.dim], v)
	Normalize(a.data[start : start+a.dim])
	return nil
}

// Grow reallocates to at least capacity, preserving contents. The only allocating
// operation on the hot path, and callers are expected to size the arena up front.
func (a *Arena) Grow(capacity int) {
	if capacity <= a.Cap() {
		return
	}
	targetCap := capacity * a.dim
	newData := make([]float32, len(a.data), targetCap)
	copy(newData, a.data)
	a.data = newData
}
