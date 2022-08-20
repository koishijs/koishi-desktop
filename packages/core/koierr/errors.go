package koierr

import "fmt"

var (
	NewErrSuccess        = func(err error) *KoiError { return NewKoiError(0, "success", err) }
	NewErrUnknown        = func(err error) *KoiError { return NewKoiError(1, "unknown error", err) }
	NewErrBadRequest     = func(err error) *KoiError { return NewKoiError(400, "bad request", err) }
	NewErrInternalError  = func(err error) *KoiError { return NewKoiError(500, "internal error", err) }
	NewErrNotImplemented = func(err error) *KoiError { return NewKoiError(501, "not implemented", err) }

	ErrSuccess        = NewErrSuccess(nil)
	ErrUnknown        = NewErrUnknown(nil)
	ErrBadRequest     = NewErrBadRequest(nil)
	ErrInternalError  = NewErrInternalError(nil)
	ErrNotImplemented = NewErrNotImplemented(nil)

	NewErrInstanceExists = func(name string) *KoiError { return NewKoiError(1100, fmt.Sprintf("instance %s already exists", name), nil) }
)
