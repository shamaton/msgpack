# MessagePack for Golang

[![Build Status](https://travis-ci.org/shamaton/msgpack.svg?branch=master)](https://travis-ci.org/shamaton/msgpack)

## Usage
### Installation
```sh
go get -u github.com/shamaton/msgpack
```

### How to use
#### use simply
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
	err = msgpack.Decode(&r, d)
	if err != nil {
		panic(err)
	}
}
```

## License

This library is under the MIT License.
