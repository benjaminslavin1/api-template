package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
)

type application struct {
	log *zerolog.Logger
}

type envParser struct {
	errs []error
}

type config struct {
	addr                           string
	debug                          bool
	srvIdleTimeout, srvReadTimeout time.Duration
	logger                         *zerolog.Logger
}

func main() {
	config, err := setupConfig()
	if err != nil {
		log.Fatal(err)
	}
	app := application{
		log: config.logger,
	}
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		select {
		case sig := <-done:
			app.log.Info().Msgf("System called: %+v: quitting...", sig)
			cancel()
		case <-ctx.Done():
			// app exited elsewhere
		}
	}()
	app.log.Debug().Msg("This is a debug statement")
	err = app.run(ctx, config)
	cancel()
	if err != nil {
		app.log.Fatal().Err(err).Msgf("Unexpected error returned from server")

	}
}

func (app *application) run(ctx context.Context, config config) error {
	srv := &http.Server{
		Addr:        config.addr,
		Handler:     app.routes(),
		IdleTimeout: config.srvIdleTimeout, // By default, Go enables keep-alives on all accepted connection. Keep-alive connections will be automatically closed after 1 minute of inactivity.
		ReadTimeout: config.srvReadTimeout, // If the request headers or body are still being read 5 seconds after the request is first accepted, then Go will close the underlying connection.
	}
	errorChannel := make(chan error, 1)
	go func() {
		app.log.Info().Msgf("server starting at %s", srv.Addr)
		errorChannel <- srv.ListenAndServe()
	}()

	select {
	case err := <-errorChannel:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-ctx.Done():
	}
	gracefulCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return srv.Shutdown(gracefulCtx)
}
