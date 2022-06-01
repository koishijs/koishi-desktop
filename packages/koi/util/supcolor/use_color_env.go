// Package supcolor determine whether the io supports color.
package supcolor

import (
	envUtil "koi/util/env"
	"strconv"
)

func UseColorEnv(env *[]string, mode int8) {
	// FORCE_COLOR
	envUtil.UseEnv(env, "FORCE_COLOR", strconv.Itoa(int(mode)))

	// COLORTERM
	if mode >= 3 {
		envUtil.UseEnv(env, "COLORTERM", "truecolor")
	}

	// TERM
	if mode >= 3 {
		envUtil.UseEnv(env, "TERM", "xterm-truecolor")
	} else if mode == 2 {
		envUtil.UseEnv(env, "TERM", "xterm-256color")
	} else if mode == 1 {
		envUtil.UseEnv(env, "TERM", "xterm-color")
	} else {
		envUtil.UseEnv(env, "TERM", "dumb")
	}

	// CLICOLOR
	if mode >= 1 {
		envUtil.UseEnv(env, "CLICOLOR", "1")
	}

	// TERM_PROGRAM
	envUtil.UseEnv(env, "TERM_PROGRAM", "Koi")
}
