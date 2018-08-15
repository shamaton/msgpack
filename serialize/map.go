package serialize

import (
	"fmt"
	"reflect"
)

func AsMap(v interface{}) ([]byte, error) {
	s := serializer{}

	// TODO : recover

	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
		if rv.Kind() == reflect.Ptr {
			rv = rv.Elem()
		}
	}
	size, err := s.calcSize(rv)
	if err != nil {
		return nil, err
	}

	s.d = make([]byte, size)
	last, err := s.create(rv, 0)
	if err != nil {
		return nil, err
	}
	if size != last {
		return nil, fmt.Errorf("failed serialization size=%d, lastIdx=%d", size, last)
	}
	return s.d, err
}
