package main

import (
	"context"
	"embed"
	"fmt"
	"html/template"
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
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	notify "github.com/iguanesolutions/go-systemd/v5/notify"
	watchdog "github.com/iguanesolutions/go-systemd/v5/notify/watchdog"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/journald"
	"github.com/urfave/cli/v2"
)

//go:embed templates/*.html assets/*
var f embed.FS

func main() {
	mqttServer := "tcp://localhost:1883"

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

func serve(log zerolog.Logger, ctx context.Context, source *source) error {
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

	var opts []graceful.Option
	for _, l := range listeners {
		opts = append(opts, graceful.WithListener(l))
	}

	r, err := graceful.New(gin.New(), opts...)
	if err != nil {
		return fmt.Errorf("failure to create gin graceful server: %w", err)
	}

	r.Use(logger.SetLogger())
	mt := multitemplate.NewRenderer()
	mt.AddFromFilesFuncs("index", htmlFunc, "templates/base.html", "templates/index.html", "templates/entry.html")
	mt.AddFromFilesFuncs("entry", htmlFunc, "templates/refresh-entry.html", "templates/entry.html")
	mt.AddFromFilesFuncs("entries", htmlFunc, "templates/entries.html", "templates/index.html", "templates/entry.html")
	mt.AddFromFiles("about", "templates/base.html", "templates/about.html")
	r.HTMLRender = mt

	r.StaticFS("/public", http.FS(f))

	r.GET("favicon.ico", func(c *gin.Context) {
		data, err := f.ReadFile("assets/mqtt-icon.svg")
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}

		c.Data(http.StatusOK, "image/svg+xml", data)
	})
	r.GET("/about", func(c *gin.Context) {
		c.HTML(http.StatusOK, "about", gin.H{})
	})
	r.GET("/", func(c *gin.Context) {
		render(c, "index", source.rpis)
	})
	r.GET("/entries", func(c *gin.Context) {
		render(c, "entries", source.rpis)
	})
	r.POST("/entry/:id/boost", func(c *gin.Context) {
		id := c.Param("id")

		log.Info().Str("Id", id).Msg("Boosting")

		token := source.client.Publish("/rpi-poe-fan/"+id+"/speed", 0, true, "100")
		<-token.Done()

		c.String(http.StatusOK, "Boost")
	})
	r.GET("/entry/:id/json", func(c *gin.Context) {
		id := c.Param("id")

		rpi, ok := source.rpis[id]
		if !ok {
			c.Status(http.StatusNotFound)
			return
		}

		c.JSON(http.StatusOK, rpi)
	})

	r.GET("/entry/:id", event.HeadersMiddleware(), sse.Middleware(), func(c *gin.Context) {
		id := c.Param("id")

		ch, err := sse.GetChannel(c)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		c.Stream(func(w io.Writer) bool {
			msg, ok := <-ch
			if !ok {
				return false
			}
			if msg == id {
				return source.emit(log, id, c)
			}
			return true
		})
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

func systemdWatchdog(log zerolog.Logger, ctx context.Context) {
	go func() {
		w, err := watchdog.New()
		if err != nil {
			log.Info().Err(err).Msg("Failure to create systemd watchdog")
			return
		}

		for {
			select {
			case <-w.NewTicker().C:
				log.Info().Msg("Watchdog ping")
				if err := w.SendHeartbeat(); err != nil {
					log.Info().Err(err).Msg("Failure to send systemd watchdog heartbeat")
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}

func systemdStopping(log zerolog.Logger, ctx context.Context) {
	go func() {
		<-ctx.Done()
		log.Info().Msg("Notifying systemd that service is stopping")
		if err := notify.Stopping(); err != nil {
			log.Info().Err(err).Msg("Failure to notify systemd that service is stopping")
		}
	}()
}

var htmlFunc = template.FuncMap{
	"pass": func(elements ...any) []any { return elements },
}

func render(c *gin.Context, template string, obj any) {
	switch c.GetHeader("Content-Type") {
	case "application/json":
		c.JSON(http.StatusOK, obj)
	default:
		c.HTML(http.StatusOK, template, obj)
	}
}
