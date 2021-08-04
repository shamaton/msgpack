# MessagePack for Golang

[![Go Reference](https://pkg.go.dev/badge/github.com/shamaton/msgpack.svg)](https://pkg.go.dev/github.com/shamaton/msgpack)
![test](https://github.com/shamaton/msgpack/workflows/test/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/shamaton/msgpack)](https://goreportcard.com/report/github.com/shamaton/msgpack)
[![codecov](https://codecov.io/gh/shamaton/msgpack/branch/master/graph/badge.svg?token=9PD2JUK5V3)](https://codecov.io/gh/shamaton/msgpack)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fshamaton%2Fmsgpack.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fshamaton%2Fmsgpack?ref=badge_shield)

## Notice
If your application serializes only primitive types, array, map and struct, code generation is also recommended.
You can get the fastest performance with [msgpackgen](https://github.com/shamaton/msgpackgen).

## Features
* Supported types : primitive / array / slice / struct / map / interface{} and time.Time
* Renaming fields via `msgpack:"field_name"`
* Omitting fields via `msgpack:"-"`
* Supports extend encoder / decoder
* Can also Encoding / Decoding struct as array

This package require more than golang version **1.13**

## Installation
```sh
go get -u github.com/shamaton/msgpack
```

## Quick Start
```go
package main

import (
  "github.com/shamaton/msgpack/v2"
)

func main() {
	type Struct struct {
		String string
	}
	v := Struct{String: "msgpack"}

	d, err := msgpack.Marshal(v)
	if err != nil {
		panic(err)
	}
	r := Struct{}
	err = msgpack.Unmarshal(d, &r)
	if err != nil {
		panic(err)
	}
}
```

## Benchmark
This result made from [shamaton/msgpack_bench](https://github.com/shamaton/msgpack_bench)
### Encode
```
BenchmarkCompareEncodeShamaton-4                  843001              1271 ns/op             320 B/op          3 allocs/op
BenchmarkCompareEncodeShamatonArray-4            1000000              1128 ns/op             256 B/op          3 allocs/op
BenchmarkCompareEncodeVmihailenco-4               308065              3562 ns/op            1000 B/op         15 allocs/op
BenchmarkCompareEncodeVmihailencoArray-4          327612              3577 ns/op            1000 B/op         15 allocs/op
BenchmarkCompareEncodeUgorji-4                    690046              1708 ns/op             904 B/op          9 allocs/op
BenchmarkCompareEncodeJson-4                      452257              2800 ns/op             824 B/op         14 allocs/op
BenchmarkCompareEncodeGob-4                       135574              8493 ns/op            2760 B/op         50 allocs/op
```

### Decode
```
BenchmarkCompareDecodeShamaton-4                  826440              1419 ns/op             512 B/op          6 allocs/op
BenchmarkCompareDecodeShamatonArray-4            1246941              1126 ns/op             512 B/op          6 allocs/op
BenchmarkCompareDecodeVmihailenco-4               234380              4914 ns/op            1055 B/op         33 allocs/op
BenchmarkCompareDecodeVmihailencoArray-4          291438              4148 ns/op             992 B/op         22 allocs/op
BenchmarkCompareDecodeUgorji-4                    497340              2302 ns/op             890 B/op         10 allocs/op
BenchmarkCompareDecodeJson-4                      152253              7632 ns/op            1144 B/op         33 allocs/op
BenchmarkCompareDecodeGob-4                        36434             34308 ns/op           10108 B/op        275 allocs/op
```


## License

This library is under the MIT License.
