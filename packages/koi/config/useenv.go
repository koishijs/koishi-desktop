package config

import "strings"

func UseConfigEnv(env []string) []string {
	if Config.Env == nil {
		return env
	}

	for _, eAddSrc := range Config.Env {
		spl := strings.Split(eAddSrc, "=")
		if len(spl) == 0 {
			continue
		}
		var eAdd string
		eKey := spl[0] + "="
		if len(spl) == 1 {
			eAdd = eKey
		} else {
			eAdd = eAddSrc
		}

		for {
			notFound := true
			for i, e := range env {
				if strings.HasPrefix(e, eKey) {
					env = append(env[:i], env[i+1:]...)
					notFound = false
					break
				}
			}

			if notFound {
				break
			}
		}

		env = append(env, eAdd)
	}

	return env
}
