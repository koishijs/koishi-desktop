package proto

import (
	"gopkg.ilharper.com/koi/core/koierr"
)

const (
	TypeResponseResult = "result"
)

type Result struct {
	// Code is the status of [proto.Result].
	// 0 represents success and any other code represents an error.
	Code uint16 `json:"code" mapstructure:"code"`

	// Data is the [proto.Result] data.
	Data any `json:"data" mapstructure:"data"`
}

func NewResult(code uint16, data any) *Response {
	return NewResponse(TypeResponseResult, &Result{
		Code: code,
		Data: data,
	})
}

func NewSuccessResult(data any) *Response {
	return NewResult(koierr.ErrSuccess.Code, data)
}

func NewErrorResult(err *koierr.KoiError) *Response {
	return NewResult(err.Code, err.Error())
}
