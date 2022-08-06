package proto

type Result struct {
	// Code is the status of [proto.Result].
	// 0 represents success and any other code represents an error.
	Code uint16 `json:"code" mapstructure:"code"`

	// Data is the [proto.Result] data.
	Data any `json:"data" mapstructure:"data"`
}

func NewResult(code uint16, data any) *Response {
	return NewResponse("result", &Result{
		Code: code,
		Data: data,
	})
}
