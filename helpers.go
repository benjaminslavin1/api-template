package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/rs/zerolog"
)

func setupConfig() (config, error) {
	var parser envParser
	isDebug := parser.bool("DEBUG", false)
	disableJSONLogs := parser.bool("DISABLE_JSON_LOGS", false)
	log := newLogger(isDebug, disableJSONLogs)
	masterConfig := config{
		addr:           parser.string("PORT", ":8080"),
		debug:          isDebug,
		srvIdleTimeout: time.Duration(parser.int("SRV_IDLE_TIMEOUT", 1)) * time.Minute,
		srvReadTimeout: time.Duration(parser.int("SRV_READ_TIMEOUT", 5)) * time.Second,
		logger:         log,
	}
	if len(parser.errs) > 0 {
		for _, err := range parser.errs {
			log.Err(err)
		}
		return config{}, errors.New("failed parsing environment variables")
	}
	return masterConfig, nil
}

func newLogger(isDebug, disableJSONLogs bool) *zerolog.Logger {
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	logLevel := zerolog.DebugLevel
	// rely on logging solution to remove debug statements
	zerolog.SetGlobalLevel(logLevel)
	if disableJSONLogs {
		logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	}
	return &logger
}

func (e *envParser) string(name, defaultVal string) string {
	val, ok := os.LookupEnv(name)
	if !ok {
		return defaultVal
	}
	return val
}

func (e *envParser) int(name string, defaultVal int) int {
	sVal, ok := os.LookupEnv(name)
	if !ok {
		return defaultVal
	}
	val, err := strconv.Atoi(sVal)
	if err != nil {
		e.errs = append(e.errs, fmt.Errorf("%v: %w", name, err))
	}
	return val
}

func (e *envParser) bool(name string, defaultVal bool) bool {
	sVal, ok := os.LookupEnv(name)
	if !ok {
		return defaultVal
	}
	val, err := strconv.ParseBool(sVal)
	if err != nil {
		e.errs = append(e.errs, fmt.Errorf("%v: %w", name, err))
	}
	return val
}

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.log.Err(err).Msg(trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}
