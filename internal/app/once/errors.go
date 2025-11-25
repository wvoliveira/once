package once

import "errors"

var (
	ErrInternalServer = errors.New("internal server error")
	ErrContentExpired = errors.New("content expired")
)
