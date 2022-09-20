//nolint:wrapcheck
package koicli

import (
	"errors"
	"fmt"

	"github.com/goccy/go-json"
	"github.com/mitchellh/mapstructure"
	"github.com/samber/do"
	"github.com/urfave/cli/v2"
	"gopkg.ilharper.com/koi/core/god/proto"
	"gopkg.ilharper.com/koi/core/koicmd"
	"gopkg.ilharper.com/koi/core/koiconfig"
	"gopkg.ilharper.com/koi/core/logger"
	"gopkg.ilharper.com/koi/sdk/client"
	"gopkg.ilharper.com/koi/sdk/manage"
)

const (
	serviceCommandStart = "gopkg.ilharper.com/koi/app/koicli/command.Start"
	serviceActionStart  = "gopkg.ilharper.com/koi/app/koicli/action.Start"
)

func newStartCommand(i *do.Injector) (*cli.Command, error) {
	do.ProvideNamed(i, serviceActionStart, newStartAction)

	return &cli.Command{
		Name:      "start",
		Usage:     "Start Instances",
		ArgsUsage: "instances",
		Action:    do.MustInvokeNamed[cli.ActionFunc](i, serviceActionStart),
	}, nil
}

func newStartAction(i *do.Injector) (cli.ActionFunc, error) {
	l := do.MustInvoke[*logger.Logger](i)

	return func(c *cli.Context) error {
		var err error

		l.Debug("Trigger action: start")

		cfg, err := do.Invoke[*koiconfig.Config](i)
		if err != nil {
			return err
		}

		manager := manage.NewKoiManager(cfg.Computed.Exe, cfg.Computed.DirLock)
		conn, err := manager.Ensure()
		if err != nil {
			return err
		}

		respC, logC, err := client.Start(
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

		var resultPs koicmd.ResultPs
		err = mapstructure.Decode(result.Data, &resultPs)
		if err != nil {
			return fmt.Errorf("failed to parse result %#+v: %w", result, err)
		}

		resultPsInstanceJSON, err := json.Marshal(resultPs)
		if err != nil {
			return fmt.Errorf("failed to marshal result %#+v: %w", resultPs, err)
		}

		fmt.Println(string(resultPsInstanceJSON))

		err = logger.Wait(respC)
		if err != nil {
			return fmt.Errorf("failed to process command ps: %w", err)
		}

		return logger.Wait(respC)
	}, nil
}
