package web

import (
	"context"
	"embed"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/Bluebugs/rpi-poe-fan/pkg/event"

	"github.com/coreos/go-systemd/activation"
	"github.com/gin-contrib/graceful"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

//go:embed assets/*
var f embed.FS

func Serve(log *zerolog.Logger, ctx context.Context, source *source, ready func() error, options ...graceful.Option) error {
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
	r.Use(APIEndpoints)

	InstantiateTemplate(r.Engine)
	ServeStaticFile(r.Engine)
	ServeDynamicPage(log, r.Engine, source, sse)
	ServeAPI(log, r.Engine, source)

	r.GET("/about", func(c *gin.Context) {
		c.HTML(http.StatusOK, "about", gin.H{})
	})

	if err := ready(); err != nil {
		return err
	}

	if err := r.RunWithContext(ctx); err != nil && err != context.Canceled {
		return fmt.Errorf("failure to run gin server: %w", err)
	}
	return nil
}
