package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"go.uber.org/fx"
)

var defaultPort = "8080"

func isRelease() bool {
	return os.Getenv("GIN_MODE") == "release"
}

func getPort() string {
	var envPort = os.Getenv("PORT")
	if envPort != "" {
		return fmt.Sprintf(":%s", envPort)
	}
	return fmt.Sprintf(":%s", defaultPort)
}

func Server(lc fx.Lifecycle, l *zerolog.Logger) *gin.Engine {

	if isRelease() {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: func(param gin.LogFormatterParams) string {
			var statusColor, methodColor, resetColor string
			if param.IsOutputColor() {
				statusColor = param.StatusCodeColor()
				methodColor = param.MethodColor()
				resetColor = param.ResetColor()
			}

			if param.Latency > time.Minute {
				param.Latency = param.Latency.Truncate(time.Second)
			}
			return fmt.Sprintf("[GIN] %v |%s %3d %s| %13v | %15s |%s %-7s %s %#v\n%s",
				param.TimeStamp.Format("2006/01/02 - 15:04:05"),
				statusColor, param.StatusCode, resetColor,
				param.Latency,
				param.ClientIP,
				methodColor, param.Method, resetColor,
				param.Path,
				param.ErrorMessage,
			)
		},
	}))
	srv := &http.Server{Addr: getPort(), Handler: router} // define a web server

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr) // the web server starts listening on 8080
			if err != nil {
				l.Error().Msgf("[My Demo] Failed to start HTTP Server at %s", srv.Addr)
				return err
			}
			go srv.Serve(ln) // process an incoming request in a go routine

			l.Info().Msgf("[My Demo] Start HTTP Server at %s", srv.Addr)
			return nil

		},
		OnStop: func(ctx context.Context) error {
			srv.Shutdown(ctx) // stop the web server
			l.Info().Msg("[My Demo] HTTP Server is stopped")
			return nil
		},
	})

	return router
}

var Module = fx.Provide(Server)
var Run = func(*gin.Engine) {}
