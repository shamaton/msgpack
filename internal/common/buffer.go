package common

import "sync"

type Buffer struct {
	Data []byte
	B1   []byte
	B2   []byte
	B4   []byte
	B8   []byte
	B16  []byte
}

var bufPool = sync.Pool{
	New: func() interface{} {
		data := make([]byte, 32)
		return &Buffer{
			Data: data,
			B1:   data[:1],
			B2:   data[:2],
			B4:   data[:4],
			B8:   data[:8],
			B16:  data[:16],
		}
	},
}

func GetBuffer() *Buffer {
	return bufPool.Get().(*Buffer)
}

func PutBuffer(buf *Buffer) {
	bufPool.Put(buf)
}

// todo : Limit
// todo : Data bytes のサイズ拡張メソッド
