package command

import (
	"code.cloudfoundry.org/lager"
	"github.com/go-kit/kit/log/level"
	"github.com/gorilla/mux"
	"github.com/kolide/kit/httputil"
	"github.com/kolide/kit/logutil"
	"github.com/kolide/kit/tlsutil"
	"github.com/lessor/lessor/pkg/broker"
	"github.com/pivotal-cf/brokerapi"
	"github.com/urfave/cli"
)

// RunBroker is the implementation of the lessor run broker command
func RunBroker() cli.Command {
	var (
		flAddr    string
		flTls     bool
		flTlsCert string
		flTlsKey  string
		flDebug   bool
	)
	return cli.Command{
		Name:  "broker",
		Usage: "Run the Open Service Broker API",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "addr",
				Value:       ":8080",
				EnvVar:      "ADDR",
				Destination: &flAddr,
				Usage:       "The address to run the server on",
			},
			cli.BoolFlag{
				Name:        "tls",
				EnvVar:      "TLS",
				Destination: &flTls,
				Usage:       "Whether or not to terminate TLS",
			},
			cli.StringFlag{
				Name:        "tls-cert",
				Value:       "",
				EnvVar:      "TLS_CERT",
				Destination: &flTlsCert,
				Usage:       "The TLS certificate",
			},
			cli.StringFlag{
				Name:        "tls-key",
				Value:       "",
				EnvVar:      "TLS_KEY",
				Destination: &flTlsKey,
				Usage:       "The TLS key",
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

			serviceBroker := &broker.Broker{
				InstanceCreators: map[string]broker.InstanceCreator{},
				InstanceBinders:  map[string]broker.InstanceBinder{},
			}

			router := mux.NewRouter()
			brokerLogger := lager.NewLogger("lessor-broker")
			brokerapi.AttachRoutes(router, serviceBroker, brokerLogger)

			srv := httputil.NewServer(flAddr, router)

			if !flTls {
				level.Info(logger).Log("msg", "starting broker", "address", flAddr, "tls", false)
				return srv.ListenAndServe()
			} else {
				level.Info(logger).Log("msg", "starting broker", "address", flAddr, "tls", true)
				srv.TLSConfig = tlsutil.NewConfig()
				return srv.ListenAndServeTLS(flTlsCert, flTlsKey)
			}

			return nil
		},
	}
}
