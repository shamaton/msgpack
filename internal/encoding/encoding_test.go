package encoding

import (
	"math"
	"reflect"
	"strconv"
	"testing"

	"github.com/shamaton/msgpack/v2/def"

	tu "github.com/shamaton/msgpack/v2/internal/common/testutil"
)

func TestEncode(t *testing.T) {
	v := 1
	vv := &v

	b, err := Encode(&vv, false)
	tu.NoError(t, err)

	tu.EqualSlice(t, b, []byte{def.PositiveFixIntMin + 1})
}

func Test_encode(t *testing.T) {
	type st struct {
		V int
	}

	type testcase struct {
		value any
		code  byte
		error error
	}

	f := func(tcs []testcase, t *testing.T) {
		for _, tc := range tcs {
			rv := reflect.ValueOf(tc.value)
			t.Run(rv.Type().String(), func(t *testing.T) {
				e := encoder{}
				size, err := e.calcSize(rv)
				tu.IsError(t, err, tc.error)
				if err != nil {
					return
				}

				e.d = make([]byte, size)
				result := e.create(rv, 0)
				tu.Equal(t, result, size)
				tu.Equal(t, e.d[0], tc.code)
			})
		}
	}

	var testcases []testcase

	// slice tests
	testcases = []testcase{
		{
			value: ([]byte)(nil),
			code:  def.Nil,
		},
		{
			value: make([]byte, math.MaxUint32+1),
			error: def.ErrUnsupportedType,
		},
		{
			value: make([]int, 1),
			code:  def.FixArray + 1,
		},
		{
			value: make([]int, math.MaxUint16),
			code:  def.Array16,
		},
		// too heavy
		//{
		//	value: make([]int, math.MaxUint32),
		//	code:  def.Array32,
		//},
		//{
		//	value: make([]int, math.MaxUint32+1),
		//	error: def.ErrUnsupportedType,
		//},
		{
			value: []st{{1}},
			code:  def.FixArray + 1,
		},
		{
			value: []chan int{make(chan int)},
			error: def.ErrUnsupportedType,
		},
	}
	t.Run("slice", func(t *testing.T) {
		f(testcases, t)
	})

	// array tests
	testcases = []testcase{
		// stack frame too large (compile error)
		//{
		//	value: [math.MaxUint32 + 1]byte{},
		//	error: def.ErrUnsupportedType,
		//},
		{
			value: [1]int{},
			code:  def.FixArray + 1,
		},
		{
			value: [math.MaxUint16]int{},
			code:  def.Array16,
		},
		// stack frame too large (compile error)
		//{
		//	value: [math.MaxUint32]int{},
		//	code:  def.Array32,
		//},
		//{
		//	value: [math.MaxUint32 + 1]int{},
		//	error: def.ErrUnsupportedType,
		//},
		{
			value: [1]st{{1}},
			code:  def.FixArray + 1,
		},
		{
			value: [1]chan int{make(chan int)},
			error: def.ErrUnsupportedType,
		},
	}
	t.Run("array", func(t *testing.T) {
		f(testcases, t)
	})

	// map tests
	createMap := func(l int) map[string]int {
		m := map[string]int{}
		for i := 0; i < l; i++ {
			m[strconv.Itoa(i)] = i
		}
		return m
	}
	testcases = []testcase{
		{
			value: (map[string]int)(nil),
			code:  def.Nil,
		},
		{
			value: createMap(1),
			code:  def.FixMap + 1,
		},
		{
			value: createMap(math.MaxUint16),
			code:  def.Map16,
		},
		// too heavy
		//{
		//	value: createMap(math.MaxUint32),
		//	code:  def.Map32,
		//},
		//{
		//	value: createMap(math.MaxUint32 + 1),
		//	error: def.ErrUnsupportedType,
		//},
		{
			value: map[chan int]int{make(chan int): 1},
			error: def.ErrUnsupportedType,
		},
		{
			value: map[string]chan int{"a": make(chan int)},
			error: def.ErrUnsupportedType,
		},
	}
	t.Run("map", func(t *testing.T) {
		f(testcases, t)
	})

	type unsupport struct {
		Chan chan int
	}

	testcases = []testcase{
		{
			value: unsupport{make(chan int)},
			error: def.ErrUnsupportedType,
		},
	}
	t.Run("struct", func(t *testing.T) {
		f(testcases, t)
	})

	ch := make(chan int)
	testcases = []testcase{
		{
			value: (*int)(nil),
			code:  def.Nil,
		},
		{
			value: &ch,
			error: def.ErrUnsupportedType,
		},
		{
			value: new(int),
			code:  0,
		},
	}
	t.Run("ptr", func(t *testing.T) {
		f(testcases, t)
	})

	type inter struct {
		V any
	}
	testcases = []testcase{
		{
			value: inter{V: make(chan int)},
			error: def.ErrUnsupportedType,
		},
		{
			value: inter{V: 1},
			code:  def.FixMap + 1,
		},
	}
	t.Run("interface", func(t *testing.T) {
		f(testcases, t)
	})
}

func Test_calcLength(t *testing.T) {
	e := encoder{}

	testcases := []struct {
		name   string
		length int
		size   int
		err    error
	}{
		{
			name:   "0x0f",
			length: 0x0f,
			size:   def.Byte1,
			err:    nil,
		},
		{
			name:   "MaxUint16",
			length: math.MaxUint16,
			size:   def.Byte1 + def.Byte2,
			err:    nil,
		},
		{
			name:   "MaxUint32",
			length: math.MaxUint32,
			size:   def.Byte1 + def.Byte4,
			err:    nil,
		},
		{
			name:   "error",
			length: math.MaxUint32 + 1,
			size:   0,
			err:    def.ErrUnsupportedLength,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			size, err := e.calcLength(tc.length)
			tu.IsError(t, err, tc.err)
			tu.Equal(t, size, tc.size)
		})
	}
}
