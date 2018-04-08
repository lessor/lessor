package command

import (
	"github.com/go-kit/kit/log/level"
	"github.com/kolide/kit/logutil"
	"github.com/urfave/cli"
)

func AdoptCluster() cli.Command {
	var (
		flDebug bool
	)
	return cli.Command{
		Name:  "cluster",
		Usage: "Adopt an existing cluster",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:        "debug",
				EnvVar:      "DEBUG",
				Destination: &flDebug,
				Usage:       "Whether or not to enable debug logging",
			},
		},
		Action: func(cliCtx *cli.Context) error {
			logger := logutil.NewCLILogger(flDebug)
			level.Info(logger).Log("msg", "lessor adopt cluster")
			return nil
		},
	}
}

func EjectCluster() cli.Command {
	var (
		flDebug bool
	)
	return cli.Command{
		Name:  "cluster",
		Usage: "Remove Lessor from an existing cluster",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:        "debug",
				EnvVar:      "DEBUG",
				Destination: &flDebug,
				Usage:       "Whether or not to enable debug logging",
			},
		},
		Action: func(cliCtx *cli.Context) error {
			logger := logutil.NewCLILogger(flDebug)
			level.Info(logger).Log("msg", "lessor eject cluster")
			return nil
		},
	}
}
