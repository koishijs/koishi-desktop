package config

import (
	envUtil "koi/util/env"
	"strings"
)

func UseConfigEnv(env *[]string) {
	if Config.Env == nil {
		return
	}

	for _, e := range Config.Env {
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

		envUtil.UseEnv(env, k, v)
	}
}
