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

	env := os.Environ()

	if cfg.Data.Isolate != "none" {
		envutil.UseEnv(&env, "HOME", cfg.Computed.DirHome)
		envutil.UseEnv(&env, "USERPROFILE", cfg.Computed.DirHome)

		if runtime.GOOS == "windows" {
			localPath := filepath.Join(cfg.Computed.DirHome, "Appdata", "Local")
			envutil.UseEnv(&env, "LOCALAPPDATA", localPath)
			roamingPath := filepath.Join(cfg.Computed.DirHome, "Appdata", "Roaming")
			envutil.UseEnv(&env, "APPDATA", roamingPath)
		}

		envutil.UseEnv(&env, "TMPDIR", cfg.Computed.DirTemp)
		envutil.UseEnv(&env, "TEMP", cfg.Computed.DirTemp)
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
	envutil.UseColorEnv(&env)
	koiconfig.UseConfigEnv(&env, cfg)

	return &env
}
