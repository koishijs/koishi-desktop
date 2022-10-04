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
	serviceCommandYarn = "gopkg.ilharper.com/koi/app/koicli/command.Yarn"
	serviceActionYarn  = "gopkg.ilharper.com/koi/app/koicli/action.Yarn"
)

func newYarnCommand(i *do.Injector) (*cli.Command, error) {
	do.ProvideNamed(i, serviceActionYarn, newYarnAction)

	return &cli.Command{
		Name:      "yarn",
		Usage:     "Run Yarn Command on Instance",
		ArgsUsage: "args",
		Action:    do.MustInvokeNamed[cli.ActionFunc](i, serviceActionYarn),

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "instance",
				Aliases: []string{"name", "n"},
				Usage:   "Target Instance",
			},
		},
	}, nil
}

func newYarnAction(i *do.Injector) (cli.ActionFunc, error) {
	l := do.MustInvoke[*logger.Logger](i)

	return func(c *cli.Context) error {
		var err error

		l.Debug("Trigger action: yarn")

		cfg, err := do.Invoke[*koiconfig.Config](i)
		if err != nil {
			return err
		}

		manager := manage.NewKoiManager(cfg.Computed.Exe, cfg.Computed.DirLock)
		conn, err := manager.Ensure(false)
		if err != nil {
			return err
		}

		respC, logC, err := client.Yarn(
			conn,
			c.String("instance"),
			c.Args().Slice(),
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
