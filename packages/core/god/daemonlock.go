package god

type DaemonLock struct {
	Pid  int    `json:"pid" mapstructure:"pid"`
	Addr string `json:"addr" mapstructure:"addr"`
}
