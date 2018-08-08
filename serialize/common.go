package serialize

import "github.com/shamaton/msgpack/def"

type common struct{}

func (c common) isPositiveFixInt64(v int64) bool {
	return def.PositiveFixIntMin <= v && v <= def.PositiveFixIntMax
}

func (c common) isPositiveFixUint64(v uint64) bool {
	return def.PositiveFixIntMin <= v && v <= def.PositiveFixIntMax
}

func (c common) isNegativeFixInt64(v int64) bool {
	return def.NegativeFixintMin <= v && v <= def.NegativeFixintMax
}
