package def

import "errors"

var (
	ErrNoData                 = errors.New("no data")
	ErrHasLeftOver            = errors.New("data has left over")
	ErrReceiverNotPointer     = errors.New("receiver not pointer")
	ErrNotMatchArrayElement   = errors.New("not match array element")
	ErrCanNotDecode           = errors.New("msgpack : invalid code")
	ErrCanNotSetSliceAsMapKey = errors.New("can not set slice as map key")
	ErrCanNotSetMapAsMapKey   = errors.New("can not set map as map key")

	ErrTooShortBytes         = errors.New("too short bytes")
	ErrLackDataLengthToSlice = errors.New("data length lacks to create slice")
	ErrLackDataLengthToMap   = errors.New("data length lacks to create map")
	ErrUnsupported           = errors.New("unsupported")
	ErrNotMatchLastIndex     = errors.New("not match last index")
)