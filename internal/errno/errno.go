package errno

import (
	"errors"
)

var (
	ErrLengthNotEqual = errors.New("Slice length not equal")
	ErrEnvelopeLen    = errors.New("Data size is less than envelope_len_size")
	ErrBinanceFutures = errors.New("Binacne futures deal error")
)
