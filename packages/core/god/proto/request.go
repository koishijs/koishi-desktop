package proto

type Request struct {
	Type string `json:"type" mapstructure:"type"`
	Data any    `json:"data" mapstructure:"data"`
}

func NewRequest(rType string, data any) *Request {
	return &Request{
		Type: rType,
		Data: data,
	}
}
