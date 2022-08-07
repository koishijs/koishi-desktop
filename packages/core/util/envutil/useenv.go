package envutil

import (
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strings"
)

var (
	titleCaser = cases.Title(language.AmericanEnglish)
)

func UseEnv(env *[]string, key string, value string) {
	RemoveEnv(env, key)
	*env = append(*env, fmt.Sprintf("%s=%s", key, value))
}

func RemoveEnv(env *[]string, key string) {
	removeEnvIntl(env, key)
	removeEnvIntl(env, titleCaser.String(key)) // hElLo => Hello
	removeEnvIntl(env, strings.ToUpper(key))
	removeEnvIntl(env, strings.ToLower(key))
}

func removeEnvIntl(env *[]string, key string) {
	for {
		notFound := true
		for i, e := range *env {
			if strings.HasPrefix(e, key+"=") {
				*env = append((*env)[:i], (*env)[i+1:]...)
				notFound = false
				break
			}
		}

		if notFound {
			break
		}
	}
}
