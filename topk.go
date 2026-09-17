package vector

// Match represents a candidate score pair.
type Match struct {
	ID    int32
	Score float32
}

// TopK keeps the k highest-scoring ids seen.
// Offer allocates nothing. For typical K (k <= 100), keeping an ordered insertion array
// avoids heap overhead and offers zero allocations.
type TopK struct {
	k     int
	items []Match
	len   int
}

// NewTopK allocates a TopK selector for capacity k.
func NewTopK(k int) *TopK {
	if k <= 0 {
		k = 1
	}
	return &TopK{
		k:     k,
		items: make([]Match, k),
		len:   0,
	}
}

// Reset resets the selector without reallocating.
func (t *TopK) Reset() {
	t.len = 0
}

// Offer considers candidate (id, score) for top-k membership.
// It keeps items sorted descending by score.
func (t *TopK) Offer(id int32, score float32) {
	if t.len < t.k {
		// Insert into sorted array
		idx := t.len
		t.len++
		for idx > 0 && t.items[idx-1].Score < score {
			t.items[idx] = t.items[idx-1]
			idx--
		}
		t.items[idx] = Match{ID: id, Score: score}
		return
	}

	// Array is full; compare with smallest score (at index t.k-1)
	if score <= t.items[t.k-1].Score {
		return
	}

	// Shift elements right until insertion point
	idx := t.k - 1
	for idx > 0 && t.items[idx-1].Score < score {
		t.items[idx] = t.items[idx-1]
		idx--
	}
	t.items[idx] = Match{ID: id, Score: score}
}

// Results appends the top matches in descending order into dst and returns dst.
func (t *TopK) Results(dst []Match) []Match {
	return append(dst, t.items[:t.len]...)
}
