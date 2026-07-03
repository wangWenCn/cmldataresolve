package cmldataresolve

import "errors"

var (
	ErrFieldNotFound   = errors.New("cmldataresolve: field not found")
	ErrValueMissing    = errors.New("cmldataresolve: value missing")
	ErrConvertFailed   = errors.New("cmldataresolve: convert failed")
	ErrInvalidContract = errors.New("cmldataresolve: invalid contract")
)
