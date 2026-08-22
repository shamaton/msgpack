package decoding

const (
	// maxPreallocElements bounds the capacity pre-allocated from an
	// attacker-declared element count before any element is decoded.
	maxPreallocElements = 1024

	// readChunkSize is the chunk size used to read a declared byte length
	// incrementally from the reader.
	readChunkSize = 32 * 1024
)

func initialCap(l int) int {
	if l < 1 {
		return 0
	}
	if l < maxPreallocElements {
		return l
	}
	return maxPreallocElements
}
