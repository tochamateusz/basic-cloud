package main

import (
	"github.com/rs/zerolog"
	"github.com/tochamateusz/basic/modules/logging"
	"go.uber.org/fx"
	fxevent "go.uber.org/fx/fxevent"
)

func InvokeTest(log *zerolog.Logger) {
	levelLog := log.Level(zerolog.TraceLevel)
	levelLog.Trace().Msgf("test %+v\n", log.GetLevel())
	log.Debug().Msgf("test %+v\n", log.GetLevel())
}

func main() {
	app := fx.New(
		logging.Module,
		fx.WithLogger(func(logger *zerolog.Logger) fxevent.Logger {
			return &logging.ZeroLogger{
				Logger: logger,
			}
		}),
		fx.Invoke(InvokeTest),
	)

	app.Run()
}
