package main

import (
	"fmt"
	"html/template"
	"io"
	"net"
	"net/http"

	"github.com/Bluebugs/rpi-poe-fan/pkg/event"
	"github.com/gin-contrib/graceful"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func NewEngine(listeners []net.Listener) (*graceful.Graceful, error) {
	var opts []graceful.Option
	for _, l := range listeners {
		opts = append(opts, graceful.WithListener(l))
	}

	r, err := graceful.New(gin.New(), opts...)
	if err != nil {
		return nil, fmt.Errorf("failure to create gin graceful server: %w", err)
	}

	return r, nil
}

func InstantiateTemplate(r *gin.Engine) {
	mt := multitemplate.NewRenderer()
	mt.AddFromFilesFuncs("index", htmlFunc, "templates/base.html", "templates/index.html", "templates/entry.html")
	mt.AddFromFilesFuncs("entry", htmlFunc, "templates/refresh-entry.html", "templates/entry.html")
	mt.AddFromFilesFuncs("entries", htmlFunc, "templates/entries.html", "templates/index.html", "templates/entry.html")
	mt.AddFromFiles("about", "templates/base.html", "templates/about.html")
	r.HTMLRender = mt
}

func ServeDynamicPage(log *zerolog.Logger, r *gin.Engine, source *source, sse *event.Event) {
	r.GET("/", func(c *gin.Context) {
		render(c, "index", source.rpis)
	})
	r.GET("/entries", func(c *gin.Context) {
		render(c, "entries", source.rpis)
	})
	r.POST("/entry/:id/boost", func(c *gin.Context) {
		id := c.Param("id")

		log.Info().Str("Id", id).Msg("Boosting")

		err := source.boost(id)
		if err != nil {
			c.Error(err)
		} else {
			c.String(http.StatusOK, "Boost")
		}
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
}

func ServeStaticFile(r *gin.Engine) {
	r.StaticFS("/public", http.FS(f))

	r.GET("favicon.ico", func(c *gin.Context) {
		data, err := f.ReadFile("assets/mqtt-icon.svg")
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}

		c.Data(http.StatusOK, "image/svg+xml", data)
	})
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
