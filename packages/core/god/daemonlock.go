package god

type DaemonLock struct {
	Pid  int    `json:"pid" mapstructure:"pid"`
	Host string `json:"host" mapstructure:"host"`
	Port string `json:"port" mapstructure:"port"`
}
