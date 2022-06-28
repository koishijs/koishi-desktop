package cli

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"github.com/virtuald/go-ordered-json"
	"io"
	"io/ioutil"
	"koi/config"
	"koi/daemon"
	"koi/util"
	l "koi/util/logger"
	"koi/util/strutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	instanceCommand = &cli.Command{
		Name:  "instance",
		Usage: "Manage instances",

		Subcommands: []*cli.Command{
			{
				Name:  "create",
				Usage: "Create new instance",

				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "name",
						Aliases:  []string{"n"},
						Usage:    "Name of the new instance",
						Required: true,
					},

					&cli.BoolFlag{
						Name:    "force",
						Aliases: []string{"f"},
						Usage:   "Empty target dir before creating.",
					},

					&cli.StringFlag{
						Name:    "with-packages",
						Aliases: []string{"p"},
						Usage:   "Install additional packages in instance",
					},

					&cli.StringFlag{
						Name:    "ref",
						Aliases: []string{"r"},
						Usage:   "The ref of the boilerplate to use",
					},

					&cli.StringFlag{
						Name:    "mirror",
						Aliases: []string{"m"},
						Usage:   "The GitHub mirror to use",
					},

					&cli.StringFlag{
						Name:    "template",
						Aliases: []string{"t"},
						Usage:   "The template repo to use",
					},
				},

				Action: createInstanceAction,
			},
		},
	}

	refRegexp = regexp.MustCompile("^[\\da-f]{40}$")
)

func createInstanceAction(c *cli.Context) error {
	var err error

	name := strutil.Trim(c.String("name"))
	l.Infof("Creating new instance: %s", name)

	var packages []string
	for _, p := range strings.Split(strutil.Trim(c.String("with-packages")), ",") {
		pp := strutil.Trim(p)
		if len(pp) > 0 {
			packages = append(packages, pp)
		}
	}
	if len(packages) > 0 {
		l.Info("With these packages:")
		for _, p := range packages {
			l.Infof("- %s", p)
		}
	}

	l.Debug("Checking target dir.")
	dir := filepath.Join(config.Config.InternalInstanceDir, name)
	if c.Bool("force") {
		l.Info("Emptying target dir.")
		err = os.RemoveAll(dir)
		if err != nil {
			l.Error("Failed to empty target dir:")
			l.Fatal(err)
		}
	} else {
		_, err := os.Stat(dir)
		if err == nil {
			entries, _ := os.ReadDir(dir)
			if len(entries) > 0 {
				l.Fatal("Instance already exists. Use '--force' if you want to recreate.")
			}
		}
	}
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		l.Error("Failed to create target dir:")
		l.Fatal(err)
	}

	l.Debug("Constructing boilerplate url.")
	mirror := strutil.Trim(c.String("mirror"))
	if mirror == "" {
		mirror = "https://github.com"
	}
	template := strutil.Trim(c.String("template"))
	if template == "" {
		template = "koishijs/boilerplate"
	}
	ref := strutil.Trim(c.String("ref"))
	if ref == "" {
		ref = "refs/heads/master"
	}
	if (!strings.HasPrefix(ref, "refs/")) && (!refRegexp.MatchString(ref)) {
		ref = "refs/heads/" + ref
	}
	boilerUrl := fmt.Sprintf("%s/%s/archive/%s.tar.gz", mirror, template, ref)

	l.Info("Downloading boilerplate from:")
	l.Info(boilerUrl)
	boilerRes, err := http.Get(boilerUrl)
	if err != nil {
		l.Error("Request to download boilerplate failed:")
		l.Fatal(err)
	}
	defer func() {
		_ = boilerRes.Body.Close()
	}()

	l.Info("[1/8] Downloading and scaffolding project.")
	err = util.Unzip(boilerRes.Body, dir, false, true)
	if err != nil {
		l.Error("Failed to unzip boilerplate.")
		l.Fatal(err)
	}

	l.Info("[2/8] Writing yarn config.")
	yarnrctmpl, err := os.Open(filepath.Join(config.Config.InternalDataDir, "yarnrc.tmpl.yml"))
	if err != nil {
		l.Fatal("Failed to open yarnrc.tmpl.yml.")
	}
	defer func() {
		_ = yarnrctmpl.Close()
	}()
	yarnrc, err := os.Create(filepath.Join(dir, ".yarnrc.yml"))
	if err != nil {
		l.Fatal("Failed to create .yarnrc.yml.")
	}
	defer func() {
		_ = yarnrc.Close()
	}()
	_, err = io.Copy(yarnrc, yarnrctmpl)
	if err != nil {
		l.Fatal("Failed to copy .yarnrc.yml.")
	}
	err = yarnrctmpl.Close()
	if err != nil {
		l.Fatal("Failed to close yarnrc.tmpl.yml.")
	}
	err = yarnrc.Close()
	if err != nil {
		l.Fatal("Failed to close .yarnrc.yml.")
	}
	yarnlock, err := os.Create(filepath.Join(dir, "yarn.lock"))
	if err != nil {
		l.Fatal("Failed to create yarn.lock.")
	}
	_ = yarnlock.Close()

	l.Info("[3/8] Installing initial packages (phase 1).")
	// Phase 1:
	// And here we've got a fresh Koishi project, with correct .yarnrc.yml inside.
	// Now, we need to manually delete devDeps from package.json.
	// The alternative `yarn workspaces focus --production --all` has been deprecated.
	// It won't create yarn.lock.
	pkgjson, err := ioutil.ReadFile(filepath.Join(dir, "package.json"))
	if err != nil {
		l.Fatal("Failed to read package.json. Invalid Koishi project?")
	}
	var pkgObj json.OrderedObject
	err = json.Unmarshal(pkgjson, &pkgObj)
	if err != nil {
		l.Fatal("Failed to parse package.json. Invalid Koishi project?")
	}
	for i, m := range pkgObj {
		// Remove `devDependencies` field
		if m.Key == "devDependencies" {
			pkgObj = append(pkgObj[:i], pkgObj[i+1:]...)
			break
		}
	}
	pkgResult, err := json.MarshalIndent(pkgObj, "", "  ")
	if err != nil {
		l.Fatal("Failed to generate package.json.")
	}
	err = os.WriteFile(filepath.Join(dir, "package.json"), pkgResult, 0666) // rw-rw-rw-
	if err != nil {
		l.Fatal("Failed to write package.json.")
	}

	l.Info("[4/8] Installing initial packages (phase 2).")
	// Phase 2:
	// Then `yarn`.
	err = daemon.RunYarnCmd(
		[]string{"--no-immutable"},
		dir,
	)
	if err != nil {
		l.Error("Err when installing initial packages.")
		l.Fatal(err)
	}

	if len(packages) > 0 {
		l.Info("[5/8] Installing additional packages (phase 3).")
		// Phase 3:
		// Then `yarn add` our additional packages.
		// Here yarn will generate an unreliable dep tree. We'll deal with it later.
		err = daemon.RunYarnCmd(
			append([]string{"add"}, packages...),
			dir,
		)
		if err != nil {
			l.Error("Err when installing additional packages.")
			l.Fatal(err)
		}

		l.Info("[6/8] Deleting node_modules (phase 4).")
		// Phase 4:
		// Now delete node_modules and yarn.lock.
		err = os.RemoveAll(filepath.Join(dir, "node_modules"))
		if err != nil {
			l.Error("Err when deleting node_modules.")
			l.Fatal(err)
		}

		l.Info("[7/8] Installing all packages (phase 5).")
		// Phase 5:
		// Finally, `yarn`.
		// This will generate the final deps we want.
		err = daemon.RunYarnCmd(
			[]string{"--no-immutable"},
			dir,
		)
		if err != nil {
			l.Error("Err when installing all packages.")
			l.Fatal(err)
		}
	}

	l.Info("[8/8] Deleting yarn cache.")
	err = os.RemoveAll(filepath.Join(dir, ".yarn"))
	if err != nil {
		l.Fatal("Failed to delete yarn cache.")
	}

	l.Info("Done. Your new instance:")
	l.Info(name)
	l.Info(dir)

	return nil
}
