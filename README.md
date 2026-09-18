# vector
<img src="docs/img/badges.svg">

Vector math, shared arena, top-k selection and little-endian float32 codec for browser-native semantic search.

## Benchmarks

`BenchmarkDot_384` is the number the whole design is judged by: it is the inner loop of
every kNN query, and the master plan derives the feasibility of the phase-3 encoder from
it. Record it in **MFLOPS**, not just ns/op — ns/op alone cannot be divided into a FLOP
budget.

One `Dot` over 384 dimensions is 384 multiplies + 383 adds ≈ **768 FLOP**.

### TinyGo → WASM (the browser target, what actually matters)

```
BenchmarkDot_384    ~233 ns/op     →  768 / 233e-9  ≈  3.3 GFLOPS
```

Measured with `tinygo test -target wasm -bench=BenchmarkDot_384 -benchtime=2s`, three runs:
229.5, 231.4, 238.1 ns/op. Scalar, no SIMD — TinyGo's WASM SIMD support is incomplete and
this package deliberately does not rely on it (see `docs/` §2).

### Native Go (reference only)

```
goos: linux, goarch: amd64, cpu: 11th Gen Intel(R) Core(TM) i7-11800H @ 2.30GHz
BenchmarkDot_384-16          37732281        93.08 ns/op   →  ≈ 8.3 GFLOPS
BenchmarkSearch_10k_384-16       1156      1058391 ns/op
BenchmarkTopK_Offer-16      635137162         1.893 ns/op
```

TinyGo/WASM runs at roughly **2.5× slower** than native here. That ratio is the useful part:
it is what any future FLOP budget for a WASM target should be discounted by.
