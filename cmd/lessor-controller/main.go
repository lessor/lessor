package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"text/tabwriter"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	clientset "github.com/lessor/lessor/pkg/client/clientset/versioned"
	informers "github.com/lessor/lessor/pkg/client/informers/externalversions"
	"github.com/lessor/lessor/pkg/controller"
	"github.com/lessor/lessor/pkg/env"
	"github.com/pkg/errors"

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

func usage(flagset *flag.FlagSet) func() {
	return func() {
		fmt.Println(`NAME:
   lessor-controller - A Kubernetes controller for managing multi-tenant workloads

OPTIONS:`)
		w := tabwriter.NewWriter(os.Stdout, 0, 2, 2, ' ', 0)
		flagset.VisitAll(func(f *flag.Flag) {
			fmt.Fprintf(w, "   --%s\t%s (default: %q)\n", f.Name, f.Usage, f.DefValue)
		})
		w.Flush()
	}
}

func newLogger(debug bool) log.Logger {
	base := log.NewJSONLogger(log.NewSyncWriter(os.Stderr))
	base = log.With(base, "ts", log.DefaultTimestampUTC)

	// setLevelKey changes the "level" key in a Go Kit logger, allowing the user
	// to set it to something else. Useful for deploying services to GCP, as
	// stackdriver expects a "severity" key instead.
	//
	// see https://github.com/go-kit/kit/issues/503
	setLevelKey := func(logger log.Logger, key interface{}) log.Logger {
		return log.LoggerFunc(func(keyvals ...interface{}) error {
			for i := 1; i < len(keyvals); i += 2 {
				if _, ok := keyvals[i].(level.Value); ok {
					// overwriting the key without copying keyvals
					// techically violates the log.Logger contract
					// but is safe in this context because none
					// of the loggers in this program retain a reference
					// to keyvals
					keyvals[i-1] = key
					break
				}
			}
			return logger.Log(keyvals...)
		})
	}
	base = setLevelKey(base, "severity")

	base = level.NewInjector(base, level.InfoValue())

	lev := level.AllowInfo()
	if debug {
		lev = level.AllowDebug()
	}

	base = log.With(base, "caller", log.Caller(6))

	var swapLogger log.SwapLogger
	swapLogger.Swap(level.NewFilter(base, lev))

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGUSR2)
	go func() {
		for {
			<-sigChan
			if debug {
				newLogger := level.NewFilter(base, level.AllowInfo())
				swapLogger.Swap(newLogger)
			} else {
				newLogger := level.NewFilter(base, level.AllowDebug())
				swapLogger.Swap(newLogger)
			}
			level.Info(&swapLogger).Log("msg", "swapping level", "debug", !debug)
			debug = !debug
		}
	}()
	return &swapLogger
}

func main() {
	if err := runController(os.Args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runController(args []string) error {
	flagset := flag.NewFlagSet("lessor-controller", flag.ExitOnError)
	var (
		flKubeConfig      = flagset.String("kubeconfig", env.String("KUBECONFIG", ""), "Path to a Kubernetes config")
		flMaster          = flagset.String("master", env.String("MASTER", ""), "Override the Kubernetes API server")
		flLocal           = flagset.Bool("local", env.Bool("LOCAL", false), "Use the local kubeconfig at ~/.kube/config")
		flWorkers         = flagset.Int("workers", env.Int("WORKERS", 2), "The number of controller workers to launch")
		flResyncPeriod    = flagset.Int("resync-period", env.Int("RESYNC_PERIOD", 60*60 /* 1 hour */), "How often to resync informers (in seconds)")
		flBroadcastEvents = flagset.Bool("broadcast-events", env.Bool("BROADCAST_EVENTS", false), "Whether or not to log events from the Kubernetes event broadcaster")
		flDebug           = flagset.Bool("debug", env.Bool("DEBUG", false), "Whether or not to enablee debug logging")
	)
	flagset.Usage = usage(flagset)
	if err := flagset.Parse(args); err != nil {
		return err
	}

	logger := newLogger(*flDebug)

	// if --local is set, use ~/.kube/config as the kubeconfig path
	kubeconfig := *flKubeConfig
	if *flLocal && kubeconfig == "" {
		kubeconfig = filepath.Join(env.String("HOME", "~/"), ".kube/config")
	}

	// get a k8s.io/client-go/rest.Config with the provided kubeconfig flags
	cfg, err := clientcmd.BuildConfigFromFlags(*flMaster, kubeconfig)
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
	// a versioned API client for the lessor.io provided types as well
	lessorClient, err := clientset.NewForConfig(cfg)
	if err != nil {
		return errors.Wrap(err, "error building clientset")
	}

	resyncPeriod := time.Duration(*flResyncPeriod) * time.Second

	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeClient, resyncPeriod)
	lessorInformerFactory := informers.NewSharedInformerFactory(lessorClient, resyncPeriod)

	c := controller.NewController(
		logger,
		kubeClient,
		lessorClient,
		kubeInformerFactory,
		lessorInformerFactory,
		*flBroadcastEvents,
	)

	stopCh := signals.SetupSignalHandler()

	go kubeInformerFactory.Start(stopCh)
	go lessorInformerFactory.Start(stopCh)

	if err = c.Run(*flWorkers, stopCh); err != nil {
		return errors.Wrap(err, "error running controller")
	}

	return nil
}
