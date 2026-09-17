package vector

import (
	"math/rand"
	"sort"
	"testing"
)

func TestSearch_MatchesNaiveReference(t *testing.T) {
	dim := 384
	count := 10000
	k := 10
	rnd := rand.New(rand.NewSource(99))

	arena := NewArena(dim, count)
	v := make([]float32, dim)

	for i := 0; i < count; i++ {
		for j := range v {
			v[j] = rnd.Float32()*2 - 1
		}
		arena.Append(v)
	}

	query := make([]float32, dim)
	for j := range query {
		query[j] = rnd.Float32()*2 - 1
	}
	Normalize(query)

	// Search using Arena.Search
	topk := NewTopK(k)
	arena.Search(query, nil, topk)
	got := topk.Results(nil)

	// Naive reference calculation
	type candidate struct {
		id    int32
		score float32
	}
	all := make([]candidate, count)
	for i := 0; i < count; i++ {
		candVec := arena.At(i)
		all[i] = candidate{
			id:    int32(i),
			score: Dot(query, candVec),
		}
	}
	sort.Slice(all, func(i, j int) bool {
		return all[i].score > all[j].score
	})
	want := all[:k]

	if len(got) != k {
		t.Fatalf("Results len = %d, want %d", len(got), k)
	}

	for i := 0; i < k; i++ {
		if got[i].ID != want[i].id {
			t.Errorf("At %d: got ID %d (score %v), want ID %d (score %v)",
				i, got[i].ID, got[i].Score, want[i].id, want[i].score)
		}
		if abs(got[i].Score-want[i].score) > 1e-5 {
			t.Errorf("At %d: score mismatch got %v, want %v", i, got[i].Score, want[i].score)
		}
	}
}

func TestSearch_KeepFilters(t *testing.T) {
	arena := NewArena(3, 5)
	for i := 0; i < 5; i++ {
		arena.Append([]float32{float32(i + 1), 0, 0})
	}

	query := []float32{1, 0, 0}
	topk := NewTopK(5)

	// Filter out even slot indices
	keep := func(i int) bool {
		return i%2 != 0
	}

	arena.Search(query, keep, topk)
	res := topk.Results(nil)

	if len(res) != 2 {
		t.Fatalf("Results len = %d, want 2", len(res))
	}
	for _, match := range res {
		if match.ID%2 == 0 {
			t.Errorf("Filtered out ID %d appeared in search results", match.ID)
		}
	}
}

func TestSearch_ZeroAllocs(t *testing.T) {
	dim := 384
	count := 1000
	rnd := rand.New(rand.NewSource(42))

	arena := NewArena(dim, count)
	v := make([]float32, dim)
	for i := 0; i < count; i++ {
		for j := range v {
			v[j] = rnd.Float32()
		}
		arena.Append(v)
	}

	query := make([]float32, dim)
	for j := range query {
		query[j] = rnd.Float32()
	}
	Normalize(query)

	topk := NewTopK(10)
	keep := func(i int) bool {
		return i%2 == 0
	}

	allocs := testing.AllocsPerRun(100, func() {
		topk.Reset()
		arena.Search(query, keep, topk)
	})

	if allocs != 0 {
		t.Errorf("Search allocated %v objects per run, want 0", allocs)
	}
}

func BenchmarkSearch_10k_384(b *testing.B) {
	dim := 384
	count := 10000
	rnd := rand.New(rand.NewSource(777))

	arena := NewArena(dim, count)
	v := make([]float32, dim)
	for i := 0; i < count; i++ {
		for j := range v {
			v[j] = rnd.Float32()
		}
		arena.Append(v)
	}

	query := make([]float32, dim)
	for j := range query {
		query[j] = rnd.Float32()
	}
	Normalize(query)

	topk := NewTopK(10)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		topk.Reset()
		arena.Search(query, nil, topk)
	}
}
