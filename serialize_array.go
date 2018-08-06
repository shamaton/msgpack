package msgpack

import "reflect"

type serializer struct {
	d []byte
}

func (s serializer) Exec(v ...interface{}) /*([]byte, error)*/ {
	for _, holder := range v {

		rv := reflect.ValueOf(holder)
		if rv.Kind() == reflect.Ptr {
			rv = rv.Elem()
			if rv.Kind() == reflect.Ptr {
				rv = rv.Elem()
			}
		}
	}

}
