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

This package require more than golang version **1.9**

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
BenchmarkCompareEncodeShamaton-4           	 1000000	      1341 ns/op	     320 B/op	       3 allocs/op
BenchmarkCompareEncodeShamatonArray-4      	 1000000	      1183 ns/op	     256 B/op	       3 allocs/op
BenchmarkCompareEncodeVmihailenco-4        	  200000	      5271 ns/op	     968 B/op	      14 allocs/op
BenchmarkCompareEncodeVmihailencoArray-4   	  300000	      5055 ns/op	     968 B/op	      14 allocs/op
BenchmarkCompareEncodeUgorji-4             	 1000000	      1772 ns/op	     872 B/op	      10 allocs/op
BenchmarkCompareEncodeZeroformatter-4      	 1000000	      1960 ns/op	     744 B/op	      13 allocs/op
BenchmarkCompareEncodeJson-4               	  300000	      3679 ns/op	    1224 B/op	      16 allocs/op
BenchmarkCompareEncodeGob-4                	  100000	     11988 ns/op	    2824 B/op	      50 allocs/op
```

### Decode
```
BenchmarkCompareDecodeShamaton-4           	 1000000	      1501 ns/op	     512 B/op	       6 allocs/op
BenchmarkCompareDecodeShamatonArray-4      	 1000000	      1032 ns/op	     512 B/op	       6 allocs/op
BenchmarkCompareDecodeVmihailenco-4        	  200000	      5573 ns/op	    1056 B/op	      33 allocs/op
BenchmarkCompareDecodeVmihailencoArray-4   	  300000	      4438 ns/op	     992 B/op	      22 allocs/op
BenchmarkCompareDecodeUgorji-4             	  500000	      2615 ns/op	     858 B/op	      11 allocs/op
BenchmarkCompareDecodeJson-4               	  200000	      9241 ns/op	    1216 B/op	      43 allocs/op
BenchmarkCompareDecodeGob-4                	   50000	     37985 ns/op	   10172 B/op	     275 allocs/op
```


## License

This library is under the MIT License.
