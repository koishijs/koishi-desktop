package koierr

var (
	ErrSuccess = NewKoiError(0, "success", nil)
	ErrUnknown = NewKoiError(1, "unknown error", nil)
)
