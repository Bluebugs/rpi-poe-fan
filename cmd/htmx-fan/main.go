package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"time"

	"github.com/Bluebugs/rpi-poe-fan/web"

	notify "github.com/iguanesolutions/go-systemd/v5/notify"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/journald"
	"github.com/urfave/cli/v2"
)

func init() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	var writer io.Writer
	writer = zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
		w.Out = os.Stderr
		w.TimeFormat = time.RFC3339
	})
	if os.Getenv("RUN_AS_SERVICE") != "" {
		writer = journald.NewJournalDWriter()
	}
	log := zerolog.New(writer).With().Timestamp().Logger()
	zerolog.DefaultContextLogger = &log
}

func main() {
	mqttServer := "tcp://localhost:1883"

	log := zerolog.DefaultContextLogger

	app := &cli.App{
		Name:        "htmx-fan",
		Description: "A little golang gin/htmx/tailwind css/mqtt server for the rpi-fan-sensor process",
		Usage:       "htmx-fan [options]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "connect-to",
				Aliases:     []string{"c"},
				Usage:       "address of the MQTT server to connect to.",
				Value:       mqttServer,
				Destination: &mqttServer,
			},
		},
		Action: func(c *cli.Context) error {
			if c.NArg() != 0 {
				return cli.Exit("unecessary parameters specified", 1)
			}

			source, err := web.Listen(mqttServer)
			if err != nil {
				return err
			}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			return web.Serve(log, ctx, source, func() error {
				if notify.IsEnabled() {
					if err := notify.Ready(); err != nil {
						return fmt.Errorf("failure to notify systemd that service is ready: %w", err)
					}
				}

				systemdWatchdog(log, ctx)
				systemdStopping(log, ctx)
				return nil
			})
		},
	}

	if info, ok := debug.ReadBuildInfo(); !ok {
		app.Version = "could not retrieve version infomration (ensure module support is activated and build again)"
	} else {
		app.Version = info.Main.Version
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
