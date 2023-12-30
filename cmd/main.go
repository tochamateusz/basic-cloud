package main

import (
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/tochamateusz/basic/api"
	"github.com/tochamateusz/basic/modules/logging"
	"github.com/tochamateusz/basic/modules/server"
	"go.uber.org/fx"
	fxevent "go.uber.org/fx/fxevent"
)

func InvokeTest(log *zerolog.Logger) {
	levelLog := log.Level(zerolog.TraceLevel)
	levelLog.Trace().Msgf("test %+v\n", log.GetLevel())
	log.Debug().Msgf("test %+v\n", log.GetLevel())
}

func main() {
	godotenv.Load()

	app := fx.New(
		logging.Module,
		fx.WithLogger(func(logger *zerolog.Logger) fxevent.Logger {
			return &logging.ZeroLogger{
				Logger: logger,
			}
		}),
		server.Module,
		fx.Invoke(server.Run),
		fx.Invoke(api.Api),
	)

	app.Run()
}
