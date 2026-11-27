package ext

// Raw holds the type and payload of a MessagePack ext value for which no
// Decoder is registered. When decoding into interface{}, an ext value that no
// registered Decoder claims is returned as Raw instead of producing an error,
// so unknown/application-specific ext payloads can still be read (and, if
// needed, re-encoded). Decoders registered via AddExtCoder / AddExtDecoder
// still take precedence; Raw is only used as a fallback.
//
// Re-encoding a Raw value preserves Type and Data but does not guarantee the
// original wire bytes: the encoder picks the shortest ext format (fixext1/
// 2/4/8/16 or ext8/16/32) for len(Data), so a non-canonical original encoding
// (e.g. a 1-byte payload written as ext8 instead of fixext1) will not
// round-trip byte-for-byte, even though the decoded Type/Data are unchanged.
type Raw struct {
	// Type is the ext type byte (application-defined, excluding format
	// extensions such as Timestamp, which registered decoders handle).
	Type int8
	// Data is the raw ext payload, excluding the ext header (format byte,
	// optional length, and type byte).
	Data []byte
}
