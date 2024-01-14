package web

import (
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/Bluebugs/rpi-poe-fan/pkg/event"
	"github.com/Bluebugs/rpi-poe-fan/web/templates"
	"github.com/gin-contrib/graceful"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func NewEngine(listeners []net.Listener, options ...graceful.Option) (*graceful.Graceful, error) {
	var opts []graceful.Option = options
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
	r.HTMLRender = &TemplRender{}
}

func ServeDynamicPage(log *zerolog.Logger, r *gin.Engine, source *source, sse *event.Event) {
	r.GET("/", func(c *gin.Context) {
		log.Info().Int("Entries", len(source.rpis)).Msg("Rendering index page")
		c.HTML(http.StatusOK, "", templates.Index(source.rpis))
	})
	r.GET("/entries", func(c *gin.Context) {
		c.HTML(http.StatusOK, "", templates.Entries(source.rpis))
	})
	r.POST("/entry/:id/boost", func(c *gin.Context) {
		id := c.Param("id")

		log.Info().Str("Id", id).Msg("Boosting")

		err := source.boost(id)
		if err != nil {
			_ = c.Error(err)
		} else {
			c.String(http.StatusOK, "Boost")
		}
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

func ServeAPI(log *zerolog.Logger, r *gin.Engine, source *source) {
	r.GET("/api/entries", func(c *gin.Context) {
		c.JSON(http.StatusOK, source.rpis)
	})
	r.GET("/api/entry/:id", func(c *gin.Context) {
		id := c.Param("id")

		rpi, ok := source.rpis[id]
		if !ok {
			c.Status(http.StatusNotFound)
			return
		}

		c.JSON(http.StatusOK, rpi)
	})
	r.POST("/api/entry/:id/boost", func(c *gin.Context) {
		id := c.Param("id")

		log.Info().Str("Id", id).Msg("Boosting")

		err := source.boost(id)
		if err != nil {
			_ = c.Error(err)
		} else {
			c.String(http.StatusOK, "Boost")
		}
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
