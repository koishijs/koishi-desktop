package koierr

var (
	NewErrSuccess        = func(err error) *KoiError { return NewKoiError(0, "success", err) }
	NewErrUnknown        = func(err error) *KoiError { return NewKoiError(1, "unknown error", err) }
	NewErrBadRequest     = func(err error) *KoiError { return NewKoiError(400, "bad request", err) }
	NewErrNotImplemented = func(err error) *KoiError { return NewKoiError(501, "not implemented", err) }

	ErrSuccess        = NewErrSuccess(nil)
	ErrUnknown        = NewErrUnknown(nil)
	ErrBadRequest     = NewErrBadRequest(nil)
	ErrNotImplemented = NewErrNotImplemented(nil)

	ErrorDict = map[uint16]*KoiError{
		0:   ErrSuccess,
		1:   ErrUnknown,
		400: ErrBadRequest,
		501: ErrNotImplemented,
	}
)
