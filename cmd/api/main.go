package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"FlashDeal-API-ELK/internal/config"
	"FlashDeal-API-ELK/internal/events"
	"FlashDeal-API-ELK/internal/httpapi"
	"FlashDeal-API-ELK/internal/logx"
	"FlashDeal-API-ELK/internal/middleware"
	"FlashDeal-API-ELK/internal/sse"
	"FlashDeal-API-ELK/internal/store"
)

func main() {
	cfg := config.Load()
	log := logx.New(cfg.ServiceName, cfg.Env)

	// choose store: types must satisfy httpapi.Store
	var st httpapi.Store

	dbPath := os.Getenv("DB_PATH")
	if dbPath != "" {
		log.Info().Str("db", dbPath).Msg("using sqlite store")
		sqlStore, err := store.NewSQL(dbPath)
		if err != nil {
			log.Fatal().Err(err).Msg("cannot init sqlite store")
		}
		st = sqlStore
	} else {
		// if you NEVER want in-memory, you can log.Fatal() here instead
		log.Info().Msg("DB_PATH empty, using in-memory store")
		st = store.NewMemStore()
	}

	// SSE hub
	hub := sse.NewHub(log)
	go hub.Run()

	// publisher (Kafka later; noop for now)
	var pub events.Publisher = events.NewNoopPublisher()

	// HTTP router
	r := httpapi.Router(httpapi.Deps{
		Log:    log,
		Store:  st,
		Hub:    hub,
		Pub:    pub,
		Config: cfg,
	})

	// HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      middleware.Chain(r, middleware.RequestID(), middleware.Recover(log), middleware.AccessLog(log)),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Info().Str("addr", srv.Addr).Msg("http server starting")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("server crashed")
		}
	}()

	// graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	log.Info().Msg("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
	hub.Close()
	log.Info().Msg("bye")
}
