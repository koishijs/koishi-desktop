package proto

const (
	TypeRequestCommand = "cmd"
)

// CommandRequest is a request of command.
// Use kebab-case for command name and flags.
type CommandRequest struct {
	// Name is the command name of this [proto.CommandRequest] in kebab-case,
	// like "config" or "dry-run".
	Name string `json:"name" mapstructure:"name"`

	// Flags is a set of flags, which use kebab-case for keys
	// and builtin types for values.
	Flags map[string]any `json:"flags" mapstructure:"flags"`
}

func NewCommandRequest(name string, flags map[string]any) *Request {
	return NewRequest(TypeRequestCommand, &CommandRequest{
		Name:  name,
		Flags: flags,
	})
}
