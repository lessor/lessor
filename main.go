package main

import (
	"math/rand"
	"path/filepath"
	"time"

	"github.com/kolide/kit/env"
	"github.com/kolide/kit/logutil"
	clientset "github.com/lessor/lessor/pkg/client/clientset/versioned"
	informers "github.com/lessor/lessor/pkg/client/informers/externalversions"
	"github.com/lessor/lessor/pkg/controller"
	"github.com/pkg/errors"
	"github.com/urfave/cli"

	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/sample-controller/pkg/signals"

	// this is required to authenticate to GKE clusters
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// runController is the implementation of the lessor run controller command
func runController() cli.Command {
	var (
		flKubeConfig      string
		flMaster          string
		flLocal           bool
		flWorkers         int
		flResyncPeriod    int
		flBroadcastEvents bool
		flDebug           bool
	)
	return cli.Command{
		Name:  "controller",
		Usage: "Run the Lessor Kubernetes Controller",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "kubeconfig",
				Value:       "",
				EnvVar:      "KUBECONFIG",
				Destination: &flKubeConfig,
				Usage:       "Path to a Kubernetes config",
			},
			cli.StringFlag{
				Name:        "master",
				Value:       "",
				EnvVar:      "MASTER",
				Destination: &flMaster,
				Usage:       "Override the Kubernetes API server",
			},
			cli.BoolFlag{
				Name:        "local",
				EnvVar:      "LOCAL",
				Destination: &flLocal,
				Usage:       "Use the local kubeconfig at ~/.kube/config",
			},
			cli.IntFlag{
				Name:        "workers",
				Value:       2,
				EnvVar:      "WORKERS",
				Destination: &flWorkers,
				Usage:       "The number of controller workers to launch",
			},
			cli.IntFlag{
				Name:        "resync-period",
				Value:       60 * 60, // 1 hour
				EnvVar:      "RESYNC_PERIOD",
				Destination: &flResyncPeriod,
				Usage:       "How often to resync informers (in seconds)",
			},
			cli.BoolFlag{
				Name:        "broadcast-events",
				EnvVar:      "BROADCAST_EVENTS",
				Destination: &flBroadcastEvents,
				Usage:       "Whether or not to log event from the Kubernetes event broadcaster",
			},
			cli.BoolFlag{
				Name:        "debug",
				EnvVar:      "DEBUG",
				Destination: &flDebug,
				Usage:       "Whether or not to enable debug logging",
			},
		},
		Action: func(cliCtx *cli.Context) error {
			logger := logutil.NewServerLogger(flDebug)

			// if --local is set, use ~/.kube/config as the kubeconfig path
			kubeconfig := flKubeConfig
			if flLocal && kubeconfig == "" {
				kubeconfig = filepath.Join(env.String("HOME", "~/"), ".kube/config")
			}

			// get a k8s.io/client-go/rest.Config with the provided kubeconfig flags
			cfg, err := clientcmd.BuildConfigFromFlags(flMaster, kubeconfig)
			if err != nil {
				return errors.Wrap(err, "error building kubeconfig")
			}

			// use the k8s.io/client-go/rest.Config to get a REST client which includes
			// a versioned API client for Kuberneetes types
			kubeClient, err := kubernetes.NewForConfig(cfg)
			if err != nil {
				return errors.Wrap(err, "error building kubernetes clientset")
			}

			// use the k8s.io/client-go/rest.Config to get a REST client which includes
			// a versioned API client for the kolide.com provided types as well
			lessorClient, err := clientset.NewForConfig(cfg)
			if err != nil {
				return errors.Wrap(err, "error building clientset")
			}

			resyncPeriod := time.Duration(flResyncPeriod) * time.Second

			kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeClient, resyncPeriod)
			lessorInformerFactory := informers.NewSharedInformerFactory(lessorClient, resyncPeriod)

			c := controller.NewController(
				logger,
				kubeClient,
				lessorClient,
				kubeInformerFactory,
				lessorInformerFactory,
				flBroadcastEvents,
			)

			stopCh := signals.SetupSignalHandler()

			go kubeInformerFactory.Start(stopCh)
			go lessorInformerFactory.Start(stopCh)

			if err = c.Run(flWorkers, stopCh); err != nil {
				return errors.Wrap(err, "error running controller")
			}

			return nil
		},
	}
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
				runController(),
			},
		},
	}

	app.RunAndExitOnError()
}
