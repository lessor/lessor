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
	app.Usage = "Running single-tenant apps in Kubernetes"
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
			Name:  "create",
			Usage: "Create resources",
			Subcommands: []cli.Command{
				createCluster(),
				putSecret(),
			},
		},
		cli.Command{
			Name:  "get",
			Usage: "Get and list resources",
			Subcommands: []cli.Command{
				getSecret(),
			},
		},
		cli.Command{
			Name:  "put",
			Usage: "Create or update resources",
			Subcommands: []cli.Command{
				putSecret(),
			},
		},
		cli.Command{
			Name:  "delete",
			Usage: "Delete resources",
			Subcommands: []cli.Command{
				deleteSecret(),
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
