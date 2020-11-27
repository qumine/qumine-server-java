package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/qumine/qumine-server-java/internal/api"
	"github.com/qumine/qumine-server-java/internal/grpc/server"
	"github.com/qumine/qumine-server-java/internal/server/operators"
	"github.com/qumine/qumine-server-java/internal/server/properties"
	su "github.com/qumine/qumine-server-java/internal/server/updater"
	"github.com/qumine/qumine-server-java/internal/server/whitelist"
	"github.com/qumine/qumine-server-java/internal/wrapper"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// ServerCommand is the subcommand for running in server mode
var ServerCommand = &cli.Command{
	Name:    "server",
	Aliases: []string{"s"},
	Usage:   "Start the QuMine Server",
	Action: func(c *cli.Context) error {
		properties.Configure()
		whitelist.Configure()
		operators.Configure()

		updateServer()
		updatePlugins()

		w := wrapper.NewWrapper()
		a := api.NewAPI(w)
		s := server.NewServer(w)

		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
		ctx, cancel := context.WithCancel(context.Background())
		wg := &sync.WaitGroup{}

		go w.Start(ctx, wg)
		go a.Start(ctx, wg)
		go s.Start(ctx, wg)

		<-interrupt
		logrus.Info("interrupt received, stopping")

		cancel()
		wg.Wait()
		logrus.Info("stopped")
		return nil
	},
}

func updateServer() {
	updater, err := su.NewUpdater()
	if err != nil {
		logrus.WithError(err).Fatal("Unsupported serverType")
	}
	// TODO: If jar exists continue
	if err := updater.Update(); err != nil {
		logrus.WithError(err).Fatal("Failed to update server")
	}
	updater = nil
}

func updatePlugins() {
	// TODO: Update all Plugins
}
