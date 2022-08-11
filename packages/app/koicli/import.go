package koicli

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/samber/do"
	"github.com/urfave/cli/v2"
	"gopkg.ilharper.com/koi/core/god/proto"
	"gopkg.ilharper.com/koi/core/koiconfig"
	"gopkg.ilharper.com/koi/core/koierr"
	"gopkg.ilharper.com/koi/core/logger"
	"gopkg.ilharper.com/koi/sdk/client"
	"gopkg.ilharper.com/koi/sdk/manage"
)

const (
	serviceActionImport = "gopkg.ilharper.com/koi/app/koicli/action.Import"
)

func newImportCommand(i *do.Injector) (*cli.Command, error) {
	do.ProvideNamed(i, serviceActionImport, newImportAction)

	return &cli.Command{
		Name:      "import",
		Usage:     "Import a Koishi Bundle",
		ArgsUsage: "path",
		Action:    do.MustInvokeNamed[cli.ActionFunc](i, serviceActionImport),

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "name",
				Aliases: []string{"n"},
				Usage:   "Name of the imported instance",
			},

			&cli.BoolFlag{
				Name:  "force",
				Usage: "Empty instance directory before creating",
			},
		},
	}, nil
}

func newImportAction(i *do.Injector) (cli.ActionFunc, error) {
	l := do.MustInvoke[*logger.Logger](i)

	return func(c *cli.Context) (err error) {
		cfg, err := do.Invoke[*koiconfig.Config](i)
		if err != nil {
			return
		}

		manager := manage.Manage(cfg.Computed.DirExe)
		conn, err := manager.Ensure()
		if err != nil {
			return
		}

		logC, respC, err := client.Import(
			conn,
			c.Args().First(),
			c.String("name"),
			c.Bool("force"),
		)

		logger.LogChannel(i, logC)

		var result proto.Result
		for {
			response := <-respC
			if response == nil {
				err = fmt.Errorf("failed to get result, response is nil")
				return
			}
			if response.Type == proto.TypeResponseResult {
				err = mapstructure.Decode(response.Data, &result)
				if err != nil {
					err = fmt.Errorf("failed to parse result: %w, response is: %v", err, response)
					return
				}
				break
			}
			// Ignore other type of responses
		}

		if result.Code != 0 {
			err = koierr.ErrorDict[result.Code]
		}

		return
	}, nil
}
