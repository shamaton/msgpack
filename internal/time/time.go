package time

var decodeAsLocal = false

// SetDecodedAsLocal sets the decoded time to local time.
func SetDecodedAsLocal(b bool) {
	decodeAsLocal = b
}
