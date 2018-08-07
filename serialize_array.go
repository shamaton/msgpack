package msgpack

import (
	"fmt"
	"reflect"
)

type serializer struct {
	d []byte
}

func (s *serializer) Exec(v interface{}) /*([]byte, error)*/ {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
		if rv.Kind() == reflect.Ptr {
			rv = rv.Elem()
		}
	}

	size := s.calcSize(rv)
	s.d = make([]byte, size)
	s.create(rv)
	fmt.Println(byte1, size, s)
}

func (s *serializer) calcSize(rv reflect.Value) int {
	ret := 0

	switch rv.Kind() {
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val := rv.Int()
		if positiveFixIntMin <= val && val <= positiveFixIntMax {
			fmt.Println("positive fix")
			ret = byte1
		} else if negativeFixintMin <= val && val <= negativeFixintMax {
			fmt.Println("negative fix")
			ret = byte1
		} else {
			ret = byte1
		}
	}
	return ret
}

func (s *serializer) create(rv reflect.Value) {

	switch rv.Kind() {
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val := rv.Int()
		if positiveFixIntMin <= val && val <= positiveFixIntMax {
			s.d[0] = byte(val)
		} else if negativeFixintMin <= val && val <= negativeFixintMax {
			s.d[0] = byte(val)
		}
	}
}
