package main

import (
	"context"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"

	spotifyauth "github.com/zmb3/spotify/v2/auth"

	"github.com/nure-mates/api/src/server/handlers"
	"github.com/nure-mates/api/src/service"

	log "github.com/sirupsen/logrus"

	"github.com/nure-mates/api/src/config"
	"github.com/nure-mates/api/src/server/http"
	"github.com/nure-mates/api/src/storage/postgres"
)

func main() {
	// read service cfg from os env
	cfg, err := config.Read()
	if err != nil {
		panic(err)
	}

	// init heroku port
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.WithError(err).Fatal("no port os var")
	}

	cfg.HTTPConfig.Port = port

	// init logger
	initLogger(cfg.LogLevel)

	log.Info("Service starting...")

	// prepare main context
	ctx, cancel := context.WithCancel(context.Background())
	setupGracefulShutdown(cancel)

	var wg = &sync.WaitGroup{}
	db, err := postgres.New(ctx, wg, &cfg.PostgresCfg)
	if err != nil {
		log.WithError(err).Fatal("postgres connection error")
	}

	spotifyRedirectURL := os.Getenv("REDIRECT_URI")
	if spotifyRedirectURL == "" {
		log.WithError(err).Fatal("got empty spotify redirect url")
	}

	log.Infof("redirect is: %s", spotifyRedirectURL)
	auth := spotifyauth.New(spotifyauth.WithRedirectURL(spotifyRedirectURL), spotifyauth.WithScopes(spotifyauth.ScopeUserReadEmail))

	srv := service.New(
		&cfg,
		db.NewAuthRepo(),
		db.NewProfileRepo(),
		db.NewTrackRepo(),
		db.NewRoomRepo(),
		auth,
	)

	httpSrv, err := http.New(
		&cfg.HTTPConfig,
		handlers.NewAuthHandler(srv),
		handlers.NewRoomHandler(srv),
		handlers.NewTrackHandler(srv),
	)

	if err != nil {
		log.WithError(err).Fatal("http server init")
	}

	httpSrv.Run(ctx, wg)

	// wait while services work
	wg.Wait()
	log.Info("Service stopped")
}

func initLogger(logLevel string) {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stderr)

	switch strings.ToLower(logLevel) {
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "trace":
		log.SetLevel(log.TraceLevel)
	default:
		log.SetLevel(log.DebugLevel)
	}
}

func setupGracefulShutdown(stop func()) {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signalChannel
		log.Error("Got Interrupt signal")
		stop()
	}()
}
