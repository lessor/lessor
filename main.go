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
	app.Usage = "A Kubernetes Operator for managing multi-tenant workloads"
	app.Version = "0.0.0"

	app.Commands = []cli.Command{
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
