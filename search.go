package vector

// Search scores query against every slot in the arena and offers the results to out.
// keep, when non-nil, filters by slot index BEFORE scoring — this is where deleted
// rows and tag filters are excluded, so a filtered query costs less, not more.
//
// Allocates nothing. query must already be normalised.
func (a *Arena) Search(query []float32, keep func(i int) bool, out *TopK) {
	if a.dim == 0 || len(query) != a.dim || out == nil {
		return
	}

	data := a.data
	dim := a.dim
	n := a.n

	for i := 0; i < n; i++ {
		if keep != nil && !keep(i) {
			continue
		}
		start := i * dim
		candidate := data[start : start+dim]

		// Dot inlined or direct call. Dot doesn't allocate.
		score := Dot(query, candidate)
		out.Offer(int32(i), score)
	}
}
