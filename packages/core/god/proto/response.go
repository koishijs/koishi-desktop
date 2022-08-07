package proto

import "gopkg.ilharper.com/x/rpl"

const (
	TypeResponseLog = "log"
)

type Response struct {
	Type string `json:"type" mapstructure:"type"`
	Data any    `json:"data" mapstructure:"data"`
}

func NewResponse(rType string, data any) *Response {
	return &Response{
		Type: rType,
		Data: data,
	}
}

func NewLog(log *rpl.Log) *Response {
	return NewResponse(TypeResponseLog, log)
}
