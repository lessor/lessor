package main

import (
	"github.com/go-kit/kit/log/level"
	"github.com/kolide/kit/logutil"
	"github.com/urfave/cli"
)

func getSecret() cli.Command {
	var (
		flDebug bool
	)
	return cli.Command{
		Name:    "secret",
		Aliases: []string{"secrets", "s"},
		Usage:   "Get a secret",
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
			level.Info(logger).Log("msg", "lessor get secret")
			return nil
		},
	}
}

func putSecret() cli.Command {
	var (
		flDebug bool
	)
	return cli.Command{
		Name:    "secret",
		Aliases: []string{"secrets", "s"},
		Usage:   "Create or update a secret",
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
			level.Info(logger).Log("msg", "lessor put secret")
			return nil
		},
	}
}

func deleteSecret() cli.Command {
	var (
		flDebug bool
	)
	return cli.Command{
		Name:    "secret",
		Aliases: []string{"secrets", "s"},
		Usage:   "Delete a secret",
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
			level.Info(logger).Log("msg", "lessor delete secret")
			return nil
		},
	}
}
