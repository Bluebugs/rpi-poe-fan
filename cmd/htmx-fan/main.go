package main

import (
	"context"
	"embed"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/Bluebugs/rpi-poe-fan/pkg/event"

	"github.com/coreos/go-systemd/activation"
	"github.com/gin-contrib/graceful"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	notify "github.com/iguanesolutions/go-systemd/v5/notify"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/journald"
	"github.com/urfave/cli/v2"
)

//go:embed templates/*.html assets/*
var f embed.FS

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
		Action: func(ctx *cli.Context) error {
			if ctx.NArg() != 0 {
				return cli.Exit("unecessary parameters specified", 1)
			}

			source, err := listen(log, mqttServer)
			if err != nil {
				return err
			}

			return serve(log, context.Background(), source)
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

func serve(log *zerolog.Logger, ctx context.Context, source *source, options ...graceful.Option) error {
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	sse := event.New()

	if err := source.subscribe(sse); err != nil {
		return err
	}
	defer source.Close()

	listeners, err := activation.Listeners()
	if err != nil {
		return fmt.Errorf("failure to get systemd listeners: %w", err)
	}

	r, err := NewEngine(listeners, options...)
	if err != nil {
		return fmt.Errorf("failure to create gin graceful server: %w", err)
	}

	r.Use(logger.SetLogger())

	InstantiateTemplate(r.Engine)
	ServeStaticFile(r.Engine)
	ServeDynamicPage(log, r.Engine, source, sse)

	r.GET("/about", func(c *gin.Context) {
		c.HTML(http.StatusOK, "about", gin.H{})
	})

	if notify.IsEnabled() {
		if err := notify.Ready(); err != nil {
			return fmt.Errorf("failure to notify systemd that service is ready: %w", err)
		}
	}

	systemdWatchdog(log, ctx)
	systemdStopping(log, ctx)

	if err := r.RunWithContext(ctx); err != nil && err != context.Canceled {
		return fmt.Errorf("failure to run gin server: %w", err)
	}
	return nil
}
