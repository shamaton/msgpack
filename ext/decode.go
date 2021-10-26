package ext

import (
	"reflect"
)

type Decoder interface {
	Code() int8
	AsValue(data []byte, k reflect.Kind) (interface{}, error)
}
