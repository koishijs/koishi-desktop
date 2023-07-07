package proc

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/samber/do"
	"gopkg.ilharper.com/koi/core/koiconfig"
	"gopkg.ilharper.com/koi/core/util"
	"gopkg.ilharper.com/koi/core/util/envutil"
)

func environ(i *do.Injector, path string) *[]string {
	cfg := do.MustInvoke[*koiconfig.Config](i)
	appVersion := do.MustInvokeNamed[string](i, util.ServiceAppVersion)
	appBuildNumber := do.MustInvokeNamed[string](i, util.ServiceAppBuildNumber)

	env := os.Environ()

	if cfg.Data.Isolate != "none" {
		envutil.UseEnv(&env, "KOI_HOST_HOME", os.Getenv("HOME"))
		envutil.UseEnv(&env, "HOME", cfg.Computed.DirHome)
		envutil.UseEnv(&env, "KOI_HOST_USERPROFILE", os.Getenv("USERPROFILE"))
		envutil.UseEnv(&env, "USERPROFILE", cfg.Computed.DirHome)

		if runtime.GOOS == "windows" {
			localPath := filepath.Join(cfg.Computed.DirHome, "Appdata", "Local")
			envutil.UseEnv(&env, "KOI_HOST_LOCALAPPDATA", os.Getenv("LOCALAPPDATA"))
			envutil.UseEnv(&env, "LOCALAPPDATA", localPath)
			roamingPath := filepath.Join(cfg.Computed.DirHome, "Appdata", "Roaming")
			envutil.UseEnv(&env, "KOI_HOST_APPDATA", os.Getenv("APPDATA"))
			envutil.UseEnv(&env, "APPDATA", roamingPath)
		}

		envutil.UseEnv(&env, "KOI_HOST_TMPDIR", os.Getenv("TMPDIR"))
		envutil.UseEnv(&env, "TMPDIR", cfg.Computed.DirTemp)
		envutil.UseEnv(&env, "KOI_HOST_TEMP", os.Getenv("TEMP"))
		envutil.UseEnv(&env, "TEMP", cfg.Computed.DirTemp)
		envutil.UseEnv(&env, "KOI_HOST_TMP", os.Getenv("TMP"))
		envutil.UseEnv(&env, "TMP", cfg.Computed.DirTemp)
	}

	// Replace PATH
	pathEnv := ""
	for {
		notFound := true
		for i, e := range env {
			if strings.HasPrefix(e, "PATH=") ||
				strings.HasPrefix(e, "Path=") ||
				strings.HasPrefix(e, "path=") {
				pathEnv = e[5:]
				env = append(env[:i], env[i+1:]...)
				notFound = false

				break
			}
		}
		if notFound {
			break
		}
	}
	var pathSepr string
	if runtime.GOOS == "windows" {
		pathSepr = ";"
	} else {
		pathSepr = ":"
	}
	if (pathEnv != "") && (cfg.Data.Isolate != "full") {
		pathEnv = path + pathSepr + pathEnv
	} else {
		pathEnv = path
	}
	env = append(env, "PATH="+pathEnv)

	envutil.UseEnv(&env, "KOISHI_AGENT", fmt.Sprintf("Koishi Desktop/%s", appVersion))
	envutil.UseEnv(&env, "KOI_APP_VERSION", appVersion)
	envutil.UseEnv(&env, "KOI_APP_BUILD_NUMBER", appBuildNumber)
	envutil.UseColorEnv(&env)
	koiconfig.UseConfigEnv(&env, cfg)

	return &env
}
