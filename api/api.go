package api

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

const API = "api"
const Version = "v1"

type ApiRouter *gin.RouterGroup

func Api(r *gin.Engine, l *zerolog.Logger) ApiRouter {
	domain, _ := os.Hostname()
	l.Info().Msgf("Domain: %s", domain)
	api := r.Group(fmt.Sprintf("%s/%s/", API, Version))
	api.GET("/", func(ctx *gin.Context) {
		l.Debug().Msg("test2")
	})

	api.GET("/test2", func(ctx *gin.Context) {
		l.Error().Msg("test5")
	})
	return api
}
