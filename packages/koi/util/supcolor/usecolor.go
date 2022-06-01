// Package supcolor determine whether the io supports color.
package supcolor

import (
	"fmt"
	"strings"
)

func UseEnvironColor(env []string, mode int8) []string {
	for {
		notFound := true
		for i, e := range env {
			if strings.HasPrefix(e, "FORCE_COLOR=") ||
				strings.HasPrefix(e, "COLORTERM=") ||
				strings.HasPrefix(e, "TERM=") ||
				strings.HasPrefix(e, "CLICOLOR=") ||
				strings.HasPrefix(e, "TERM_PROGRAM=") {
				env = append(env[:i], env[i+1:]...)
				notFound = false
				break
			}
		}

		if notFound {
			break
		}
	}

	// FORCE_COLOR
	env = append(env, fmt.Sprintf("FORCE_COLOR=%d", mode))

	// COLORTERM
	if mode >= 3 {
		env = append(env, "COLORTERM=truecolor")
	}

	// TERM
	if mode >= 3 {
		env = append(env, "TERM=xterm-truecolor")
	} else if mode == 2 {
		env = append(env, "TERM=xterm-256color")
	} else if mode == 1 {
		env = append(env, "TERM=xterm-color")
	} else {
		env = append(env, "TERM=dumb")
	}

	// CLICOLOR
	if mode >= 1 {
		env = append(env, "CLICOLOR=1")
	}

	// TERM_PROGRAM
	env = append(env, "TERM_PROGRAM=Koi")

	return env
}
