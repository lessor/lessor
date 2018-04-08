package main

import (
	"math/rand"
	"time"

	"github.com/kolide/kit/version"
	"github.com/urfave/cli"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	app := cli.NewApp()
	app.Name = "lessor"
	app.Usage = "Deploy, manage, and secure single-tenant applications on Kubernetes"
	app.Version = version.Version().Version
	cli.VersionPrinter = func(c *cli.Context) {
		version.PrintFull()
	}

	app.Commands = []cli.Command{
		cli.Command{
			Name:  "adopt",
			Usage: "Adopt an existing cluster",
			Subcommands: []cli.Command{
				adoptCluster(),
			},
		},
		cli.Command{
			Name:  "run",
			Usage: "Run server workloads",
			Subcommands: []cli.Command{
				runController(),
			},
		},
	}

	app.RunAndExitOnError()
}
