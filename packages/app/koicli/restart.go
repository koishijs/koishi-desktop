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
	serviceCommandRestart = "gopkg.ilharper.com/koi/app/koicli/command.Restart"
	serviceActionRestart  = "gopkg.ilharper.com/koi/app/koicli/action.Restart"
)

func newRestartCommand(i *do.Injector) (*cli.Command, error) {
	do.ProvideNamed(i, serviceActionRestart, newRestartAction)

	return &cli.Command{
		Name:      "restart",
		Usage:     "Restart Instances",
		ArgsUsage: "instances",
		Action:    do.MustInvokeNamed[cli.ActionFunc](i, serviceActionRestart),
	}, nil
}

func newRestartAction(i *do.Injector) (cli.ActionFunc, error) {
	l := do.MustInvoke[*logger.Logger](i)

	return func(c *cli.Context) error {
		var err error

		l.Debug("Trigger action: restart")

		cfg, err := do.Invoke[*koiconfig.Config](i)
		if err != nil {
			return err
		}

		manager := manage.NewKoiManager(cfg.Computed.Exe, cfg.Computed.DirLock)
		conn, err := manager.Ensure(false)
		if err != nil {
			return err
		}

		respC, logC, err := client.Restart(
			conn,
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
					return fmt.Errorf("failed to parse response %#+v: %w", response, err)
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
