package koiconfig

import (
	"gopkg.ilharper.com/koi/core/util/envutil"
	"strings"
)

func UseConfigEnv(env *[]string, cfg *Config) {
	if cfg.Data.Env == nil {
		return
	}

	for _, e := range cfg.Data.Env {
		if len(e) == 0 {
			continue
		}

		i := strings.Index(e, "=")
		var k, v string
		if i >= 0 {
			k = e[:i]
			v = e[i+1:]
		} else {
			k = e
			v = ""
		}

		envutil.UseEnv(env, k, v)
	}
}
