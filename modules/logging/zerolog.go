package logging

import (
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"runtime"
	"strings"

	"github.com/rs/zerolog"
	fxevent "go.uber.org/fx/fxevent"
)

// ZeroLogger is an Fx event logger that logs events to Zero.
type ZeroLogger struct {
	Logger *zerolog.Logger
}

var vendorRe = regexp.MustCompile("^.*?/vendor/")

// sanitize makes the function name suitable for logging display. It removes
// url-encoded elements from the `dot.git` package names and shortens the
// vendored paths.
func sanitize(function string) string {
	// Use the stdlib to un-escape any package import paths which can happen
	// in the case of the "dot-git" postfix. Seems like a bug in stdlib =/
	if unescaped, err := url.QueryUnescape(function); err == nil {
		function = unescaped
	}

	// strip everything prior to the vendor
	return vendorRe.ReplaceAllString(function, "vendor/")
}

// FuncName returns a funcs formatted name
func FuncName(fn interface{}) string {
	fnV := reflect.ValueOf(fn)
	if fnV.Kind() != reflect.Func {
		return fmt.Sprint(fn)
	}

	function := runtime.FuncForPC(fnV.Pointer()).Name()
	return fmt.Sprintf("%s()", sanitize(function))
}

// LogEvent logs the given event to the provided Zap logger.
func (l *ZeroLogger) LogEvent(event fxevent.Event) {
	switch e := event.(type) {
	case *fxevent.OnStartExecuting:
		l.Logger.Info().
			Str("callee", e.FunctionName).
			Str("caller", e.CallerName).
			Msg("hook executing")
	case *fxevent.OnStartExecuted:
		if e.Err != nil {
			l.Logger.Error().
				Err(e.Err).
				Str("method", e.Method).
				Str("callee", e.FunctionName).
				Str("caller", e.CallerName).
				Msg("hook execute failed")
		} else {
			l.Logger.Info().
				Str("method", e.Method).
				Str("callee", e.FunctionName).
				Str("caller", e.CallerName).
				Str("runtime", e.Runtime.String()).
				Msg("hook executing")
		}
	case *fxevent.Supplied:
		l.Logger.Info().
			Str("type", e.TypeName).
			Msg("supplied")
	case *fxevent.Provided:
		for _, rtype := range e.OutputTypeNames {
			l.Logger.Info().
				Str("constructor", FuncName(e.ConstructorName)).
				Str("type", rtype).
				Msg("provided")
		}
		if e.Err != nil {
			l.Logger.Error().
				Err(e.Err).
				Msg("error encountered while applying options")
		}
	case *fxevent.Invoked:
		if e.Err != nil {
			l.Logger.Error().
				Err(e.Err).
				Str("stack", e.Trace).
				Str("function", FuncName(e.FunctionName)).
				Msg("invoke failed")
		} else {
			l.Logger.Info().
				Str("function", FuncName(e.FunctionName)).
				Msg("invoked")
		}
	case *fxevent.Stopped:
		if e.Err != nil {
			l.Logger.Error().
				Err(e.Err).
				Msg("stop failed")
		}
	case *fxevent.Stopping:
		l.Logger.Info().
			Str("signal", strings.ToUpper(e.Signal.String())).
			Msg("received signal")
	case *fxevent.RolledBack:
		l.Logger.Error().
			Err(e.Err).
			Msg("rollback failed")
	case *fxevent.RollingBack:
		l.Logger.Error().
			Err(e.StartErr).
			Msg("start failed, rolling back")
	case *fxevent.Started:
		if e.Err != nil {
			l.Logger.Error().
				Err(e.Err).
				Msg("start failed")
		} else {
			l.Logger.Info().
				Msg("started")
		}
	case *fxevent.LoggerInitialized:
		if e.Err != nil {
			l.Logger.Error().
				Err(e.Err).
				Msg("custom logger installation failed")
		} else {
			l.Logger.Info().
				Str("function", FuncName(e.ConstructorName)).
				Msg("installed custom fxevent.Logger")
		}
	}
}
