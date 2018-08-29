package deserialize

import "reflect"

// todo : move to slice or map
var (
	typeIntSlice   = reflect.TypeOf([]int{})
	typeInt8Slice  = reflect.TypeOf([]int8{})
	typeInt16Slice = reflect.TypeOf([]int16{})
	typeInt32Slice = reflect.TypeOf([]int32{})
	typeInt64Slice = reflect.TypeOf([]int64{})

	typeMapStringInt = reflect.TypeOf(map[string]int{})
)
