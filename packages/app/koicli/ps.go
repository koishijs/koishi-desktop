package koicli

import (
	"errors"
	"fmt"

	"github.com/goccy/go-json"
	"github.com/mitchellh/mapstructure"
	"github.com/samber/do"
	"github.com/urfave/cli/v2"
	"golang.org/x/text/message"
	"gopkg.ilharper.com/koi/core/god/proto"
	"gopkg.ilharper.com/koi/core/koicmd"
	"gopkg.ilharper.com/koi/core/koiconfig"
	"gopkg.ilharper.com/koi/core/logger"
	"gopkg.ilharper.com/koi/sdk/client"
	"gopkg.ilharper.com/koi/sdk/manage"
)

const (
	serviceCommandPs = "gopkg.ilharper.com/koi/app/koicli/command.Ps"
	serviceActionPs  = "gopkg.ilharper.com/koi/app/koicli/action.Ps"
)

func newPsCommand(i *do.Injector) (*cli.Command, error) {
	p := do.MustInvoke[*message.Printer](i)

	do.ProvideNamed(i, serviceActionPs, newPsAction)

	return &cli.Command{
		Name:   "ps",
		Usage:  p.Sprintf("Show Process Status"),
		Action: do.MustInvokeNamed[cli.ActionFunc](i, serviceActionPs),

		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "all",
				Aliases: []string{"a"},
				Usage:   p.Sprintf("Including stopped instances"),
			},
		},
	}, nil
}

func newPsAction(i *do.Injector) (cli.ActionFunc, error) {
	l := do.MustInvoke[*logger.Logger](i)
	p := do.MustInvoke[*message.Printer](i)

	return func(c *cli.Context) error {
		var err error

		l.Debug(p.Sprintf("Trigger action: ps"))

		cfg, err := do.Invoke[*koiconfig.Config](i)
		if err != nil {
			return err
		}

		manager := manage.NewKoiManager(cfg.Computed.Exe, cfg.Computed.DirLock)
		conn, err := manager.Ensure(true)
		if err != nil {
			return errors.New(p.Sprintf("failed to get daemon connection: %v", err))
		}

		respC, logC, err := client.Ps(conn, c.Bool("all"))
		if err != nil {
			return errors.New(p.Sprintf("failed to process command ps: %v", err))
		}

		logger.LogChannel(i, logC)

		var result proto.Result
		for {
			response := <-respC
			if response == nil {
				return errors.New(p.Sprintf("failed to get result, response is nil"))
			}
			if response.Type == proto.TypeResponseResult {
				err = mapstructure.Decode(response.Data, &result)
				if err != nil {
					return errors.New(p.Sprintf("failed to parse response %#+v: %v", response, err))
				}

				break
			}
			// Ignore other type of responses
		}

		if result.Code != 0 {
			s, ok := result.Data.(string)
			if !ok {
				return errors.New(p.Sprintf("result data %#+v is not string", result.Data))
			}

			return errors.New(s)
		}

		var resultPs koicmd.ResultPs
		err = mapstructure.Decode(result.Data, &resultPs)
		if err != nil {
			return errors.New(p.Sprintf("failed to parse result %#+v: %v", result, err))
		}

		resultPsInstanceJSON, err := json.Marshal(resultPs)
		if err != nil {
			return errors.New(p.Sprintf("failed to marshal result %#+v: %v", resultPs, err))
		}

		fmt.Println(string(resultPsInstanceJSON))

		err = logger.Wait(respC)
		if err != nil {
			return errors.New(p.Sprintf("failed to process command ps: %v", err))
		}

		return nil
	}, nil
}
