package vector

// Quantizer compresses an arena to int8 with a per-vector scale, cutting resident
// memory by 4×. Needed past roughly 100k documents (master index D1).
// NOT IMPLEMENTED — Phase 5.
type Quantizer struct{}
