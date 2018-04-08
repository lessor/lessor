package main

import (
	"math/rand"
	"time"

	"github.com/lessor/lessor/pkg/command"
	"github.com/urfave/cli"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	app := cli.NewApp()
	app.Name = "lessor"
	app.Usage = "Deploy, manage, and secure single-tenant applications on Kubernetes"
	app.Version = "0.0.0"

	app.Commands = []cli.Command{
		cli.Command{
			Name:  "adopt",
			Usage: "Adopt an existing cluster",
			Subcommands: []cli.Command{
				command.AdoptCluster(),
			},
		},
		cli.Command{
			Name:  "eject",
			Usage: "Remove Lessor from an existing cluster",
			Subcommands: []cli.Command{
				command.EjectCluster(),
			},
		},
		cli.Command{
			Name:  "run",
			Usage: "Run server workloads",
			Subcommands: []cli.Command{
				command.RunController(),
			},
		},
	}

	app.RunAndExitOnError()
}
