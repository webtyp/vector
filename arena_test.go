package vector

import "testing"

func TestArena_AppendAtRoundTrip(t *testing.T) {
	arena := NewArena(3, 2)
	v1 := []float32{3, 4, 0}
	idx1, err := arena.Append(v1)
	if err != nil {
		t.Fatalf("Append failed: %v", err)
	}
	if idx1 != 0 {
		t.Errorf("Expected index 0, got %d", idx1)
	}

	got1 := arena.At(0)
	if len(got1) != 3 {
		t.Fatalf("Expected dim 3, got %d", len(got1))
	}
	// Normalised [3, 4, 0] is [0.6, 0.8, 0]
	if abs(got1[0]-0.6) > 1e-6 || abs(got1[1]-0.8) > 1e-6 || got1[2] != 0 {
		t.Errorf("At(0) = %v, want normalized [0.6, 0.8, 0]", got1)
	}

	// Ensure v1 was not retained/mutated by callers
	if v1[0] != 3 || v1[1] != 4 {
		t.Errorf("Original vector v1 mutated: %v", v1)
	}
}

func TestArena_AtIsAView(t *testing.T) {
	arena := NewArena(2, 1)
	arena.Append([]float32{1, 0})

	view := arena.At(0)
	view[0] = 0.5

	got := arena.At(0)
	if got[0] != 0.5 {
		t.Errorf("Arena contents did not reflect write to view. Got %v, want 0.5", got[0])
	}
}

func TestArena_GrowPreserves(t *testing.T) {
	arena := NewArena(2, 1)
	arena.Append([]float32{1, 0})
	arena.Append([]float32{0, 1})

	if arena.Len() != 2 {
		t.Fatalf("Len = %d, want 2", arena.Len())
	}

	arena.Grow(10)
	if arena.Cap() < 10 {
		t.Errorf("Cap = %d, want >= 10", arena.Cap())
	}
	if arena.Len() != 2 {
		t.Errorf("Len after grow = %d, want 2", arena.Len())
	}

	got0 := arena.At(0)
	got1 := arena.At(1)
	if got0[0] != 1.0 || got1[1] != 1.0 {
		t.Errorf("Contents altered after Grow: got0=%v, got1=%v", got0, got1)
	}
}
