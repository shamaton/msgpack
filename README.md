# MessagePack for Golang

[![GoDoc](https://godoc.org/github.com/shamaton/msgpack?status.svg)](https://godoc.org/github.com/shamaton/msgpack)
[![Build Status](https://travis-ci.org/shamaton/msgpack.svg?branch=master)](https://travis-ci.org/shamaton/msgpack)
[![Coverage Status](https://coveralls.io/repos/github/shamaton/msgpack/badge.svg)](https://coveralls.io/github/shamaton/msgpack)
[![Releases](https://img.shields.io/github/release/shamaton/msgpack.svg)](https://github.com/shamaton/msgpack/releases)

* Supported types : primitive / array / slice / struct / map / interface{} and time.Time
* Renames fields via `msgpack:"field_name"`
* Ignores fields via `msgpack:"ignore"`
* Supports extend encoder / decoder
* Can also Encoding / Decoding struct as array

## Installation
```sh
go get -u github.com/shamaton/msgpack
```

## Quick Start
```go
package main;

import (
  "github.com/shamaton/msgpack"
)

func main() {
	type Struct struct {
		String string
	}
	v := Struct{String: "msgpack"}

	d, err := msgpack.Encode(v)
	if err != nil {
		panic(err)
	}
	r := Struct{}
	err = msgpack.Decode(d, &r)
	if err != nil {
		panic(err)
	}
}
```

## Benchmark
This result made from [shamaton/msgpack_bench](https://github.com/shamaton/msgpack_bench)
### Encode
```
BenchmarkCompareEncodeShamaton-8           	 1000000	      1319 ns/op	     320 B/op	       3 allocs/op
BenchmarkCompareEncodeShamatonArray-8      	 1000000	      1172 ns/op	     256 B/op	       3 allocs/op
BenchmarkCompareEncodeVmihailenco-8        	  300000	      4626 ns/op	    1000 B/op	      14 allocs/op
BenchmarkCompareEncodeVmihailencoArray-8   	  500000	      3918 ns/op	     680 B/op	      13 allocs/op
BenchmarkCompareEncodeUgorji-8             	 1000000	      1985 ns/op	     986 B/op	      11 allocs/op
BenchmarkCompareEncodeJson-8               	  500000	      3649 ns/op	    1224 B/op	      16 allocs/op
BenchmarkCompareEncodeGob-8                	  100000	     12324 ns/op	    2824 B/op	      50 allocs/op
```

### Decode
```
BenchmarkCompareDecodeShamaton-8           	 1000000	      1494 ns/op	     512 B/op	       6 allocs/op
BenchmarkCompareDecodeShamatonArray-8      	 1000000	      1031 ns/op	     512 B/op	       6 allocs/op
BenchmarkCompareDecodeVmihailenco-8        	  200000	      5660 ns/op	    1056 B/op	      33 allocs/op
BenchmarkCompareDecodeVmihailencoArray-8   	  300000	      4779 ns/op	     992 B/op	      22 allocs/op
BenchmarkCompareDecodeUgorji-8             	  500000	      2774 ns/op	     844 B/op	      12 allocs/op
BenchmarkCompareDecodeJson-8               	  200000	      9721 ns/op	    1216 B/op	      43 allocs/op
BenchmarkCompareDecodeGob-8                	   50000	     37553 ns/op	   10172 B/op	     275 allocs/op
```


## License

This library is under the MIT License.
