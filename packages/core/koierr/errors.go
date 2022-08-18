package koierr

var (
	ErrSuccess        = NewKoiError(0, "success", nil)
	ErrUnknown        = NewKoiError(1, "unknown error", nil)
	ErrNotImplemented = NewKoiError(501, "not implemented", nil)

	ErrorDict = map[uint16]*KoiError{
		0:   ErrSuccess,
		1:   ErrUnknown,
		501: ErrNotImplemented,
	}
)
