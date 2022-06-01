package supcolor

import (
	"koi/util/isatty"
	"os"
	"strconv"
	"strings"
)

var (
	Stdout int8
	Stderr int8
)

func SupColor(stream *os.File) int8 {
	env := os.Environ()

	for _, e := range env {
		if strings.HasPrefix(e, "FORCE_COLOR=") {
			// Use user forced color mode
			s := e[12:]
			if len(s) == 0 {
				return 1
			}
			if s == "true" {
				return 1
			}
			if s == "false" {
				return 0
			}
			i64, err := strconv.ParseInt(s, 10, 8)
			i := int8(i64)
			if err != nil {
				// Something that's not integer
				return 1
			}
			if i < 3 {
				return i
			}
			return 3
		}
	}

	if stream == nil {
		// There's even no stream??
		return 0
	}

	if !isatty.Isatty(stream.Fd()) {
		// Isn't a TTY
		return 0
	}

	for _, e := range env {
		if e == "TERM=dumb" {
			return 0
		}

		if strings.HasPrefix(e, "CI=") {
			if strings.HasPrefix(e, "TRAVIS=") ||
				strings.HasPrefix(e, "CIRCLECI=") ||
				strings.HasPrefix(e, "APPVEYOR=") ||
				strings.HasPrefix(e, "GITLAB_CI=") ||
				strings.HasPrefix(e, "GITHUB_ACTIONS=") ||
				strings.HasPrefix(e, "BUILDKITE=") ||
				strings.HasPrefix(e, "DRONE=") {
				return 1
			}
		}

		if e == "CI_NAME=codeship" {
			return 1
		}

		if strings.HasPrefix(e, "TEAMCITY_VERSION=") {
			return 1
		}

		if strings.HasPrefix(e, "TF_BUILD=") &&
			strings.HasPrefix(e, "AGENT_NAME=") {
			return 1
		}

		if strings.HasPrefix(e, "COLORTERM=") {
			if e == "COLORTERM=truecolor" {
				return 3
			}
			return 1
		}

		if strings.HasPrefix(e, "TERM=") {
			s := e[5:]

			if strings.Contains(s, "256") {
				return 3
			}

			if strings.Contains(s, "screen") ||
				strings.Contains(s, "xterm") ||
				strings.Contains(s, "vt100") ||
				strings.Contains(s, "vt220") ||
				strings.Contains(s, "rxvt") ||
				strings.Contains(s, "color") ||
				strings.Contains(s, "ansi") ||
				strings.Contains(s, "cygwin") ||
				strings.Contains(s, "linux") {
				return 1
			}
		}
	}

	return 0
}

func init() {
	Stdout = SupColor(os.Stdout)
	Stderr = SupColor(os.Stderr)
}
