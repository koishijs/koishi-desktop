//nolint:wrapcheck
package koicli

import (
	"errors"
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/samber/do"
	"github.com/urfave/cli/v2"
	"gopkg.ilharper.com/koi/core/god/proto"
	"gopkg.ilharper.com/koi/core/koiconfig"
	"gopkg.ilharper.com/koi/core/logger"
	"gopkg.ilharper.com/koi/sdk/client"
	"gopkg.ilharper.com/koi/sdk/manage"
)

const (
	serviceCommandImport = "gopkg.ilharper.com/koi/app/koicli/command.Import"
	serviceActionImport  = "gopkg.ilharper.com/koi/app/koicli/action.Import"
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

	return func(c *cli.Context) error {
		var err error

		l.Debug("Trigger action: import")

		cfg, err := do.Invoke[*koiconfig.Config](i)
		if err != nil {
			return err
		}

		manager := manage.NewKoiManager(cfg.Computed.Exe, cfg.Computed.DirLock)
		conn, err := manager.Ensure(false)
		if err != nil {
			return err
		}

		respC, logC, err := client.Import(
			conn,
			c.Args().First(),
			c.String("name"),
			c.Bool("force"),
		)
		if err != nil {
			return err
		}

		logger.LogChannel(i, logC)

		var result proto.Result
		for {
			response := <-respC
			if response == nil {
				return fmt.Errorf("failed to get result, response is nil")
			}
			if response.Type == proto.TypeResponseResult {
				err = mapstructure.Decode(response.Data, &result)
				if err != nil {
					return fmt.Errorf("failed to parse result %#+v: %w", response, err)
				}

				break
			}
			// Ignore other type of responses
		}

		if result.Code != 0 {
			s, ok := result.Data.(string)
			if !ok {
				return fmt.Errorf("result data %#+v is not string", result.Data)
			}

			return errors.New(s)
		}

		return logger.Wait(respC)
	}, nil
}
