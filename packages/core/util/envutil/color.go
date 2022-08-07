package envutil

func UseColorEnv(env *[]string) {
	UseEnv(env, "FORCE_COLOR", "3")
	UseEnv(env, "COLORTERM", "truecolor")
	UseEnv(env, "TERM", "xterm-truecolor")
	UseEnv(env, "CLICOLOR", "1")
	UseEnv(env, "TERM_PROGRAM", "Koi")
}
