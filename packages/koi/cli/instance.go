package cli

import (
	"github.com/urfave/cli/v2"
	"koi/config"
	"koi/daemon"
	"path"
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
)

func createInstanceAction(c *cli.Context) error {
	l.Debug("Now create instance.")

	name := strings.Trim(c.String("name"), " ")
	l.Infof("Creating new instance: %s", name)

	var packages []string
	for _, p := range strings.Split(strings.Trim(c.String("with-packages"), " "), ",") {
		pp := strings.Trim(p, " ")
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

	l.Debug("Constructing args.")
	args := []string{"init", "koishi", name, "-y"}
	if ref := strings.Trim(c.String("ref"), " "); ref != "" {
		args = append(args, "-r", ref)
	}
	if mirror := strings.Trim(c.String("mirror"), " "); mirror != "" {
		args = append(args, "-m", mirror)
	}
	if template := strings.Trim(c.String("template"), " "); template != "" {
		args = append(args, "-t", template)
	}
	l.Debug("Constructed. Args:")
	l.Debug(args)

	l.Debug("Now init koishi.")
	err := daemon.RunNodeCmd(
		"npm",
		args,
		config.Config.InternalInstanceDir,
	)
	if err != nil {
		l.Error("Err when initializing koishi.")
		l.Fatal(err)
	}

	dir := path.Join(config.Config.InternalInstanceDir, name)
	l.Debug("Now install packages.")
	args = []string{"yarn", "add", "yarn", "-W", "--production"}
	args = append(args, packages...)
	err = daemon.RunNodeCmd(
		"npx",
		args,
		dir,
	)
	if err != nil {
		l.Error("Err when installing packages.")
		l.Fatal(err)
	}

	l.Info("Done. Your new instance:")
	l.Info(name)
	l.Info(dir)

	return nil
}
