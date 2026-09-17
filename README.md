# vector
<img src="docs/img/badges.svg">

Vector math, shared arena, top-k selection and little-endian float32 codec for browser-native semantic search.

## Benchmarks

```
goos: linux
goarch: amd64
pkg: webtyp.com/vector
cpu: Intel(R) Xeon(R) Processor @ 2.30GHz
BenchmarkDot_384-4          	 6611556	       181.1 ns/op
BenchmarkSearch_10k_384-4   	     492	   2.18 ms/op
BenchmarkTopK_Offer-4       	298609605	         4.02 ns/op
```
