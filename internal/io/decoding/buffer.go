package decoding

import "sync"

type buffer struct {
	data []byte
	b1   []byte
	b2   []byte
	b4   []byte
	b8   []byte
	b16  []byte
}

var bufPool = sync.Pool{
	New: func() interface{} {
		data := make([]byte, 32)
		return &buffer{
			data: data,
			b16:  data[:16],
			b8:   data[:8],
			b4:   data[:4],
			b2:   data[:2],
			b1:   data[:1],
		}
	},
}

//buf := bufPool.Get().(*buffer)
//data := encode(buf.data) // reuse buf.data
//
//newBuf := make([]byte, len(data))
//copy(newBuf, buf)
//
//buf.data = data
//bufPool.Put(buf)
