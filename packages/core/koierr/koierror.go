package koierr

type KoiError struct {
	Code uint16
	Msg  string
	Err  error
}

func NewKoiError(code uint16, msg string, err error) *KoiError {
	return &KoiError{Code: code, Msg: msg, Err: err}
}

func (e *KoiError) Error() string {
	if e.Err != nil {
		return "koi: " + e.Msg + ": " + e.Err.Error()
	} else {
		return "koi: " + e.Msg
	}
}

func (e *KoiError) Unwrap() error {
	return e.Err
}
