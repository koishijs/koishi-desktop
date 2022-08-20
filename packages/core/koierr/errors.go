package koierr

var (
	ErrSuccess        = NewKoiError(0, "success", nil)
	ErrUnknown        = NewKoiError(1, "unknown error", nil)
	ErrBadRequest     = NewKoiError(400, "bad request", nil)
	ErrNotImplemented = NewKoiError(501, "not implemented", nil)

	ErrorDict = map[uint16]*KoiError{
		0:   ErrSuccess,
		1:   ErrUnknown,
		400: ErrBadRequest,
		501: ErrNotImplemented,
	}
)
