package httpapi

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
)

// Router wires the HTTP routes using Deps from handlers.go.
func Router(d Deps) http.Handler {
	r := chi.NewRouter()

	// Normalize paths like `/deals/` -> `/deals`
	r.Use(chimw.RedirectSlashes)

	// Health
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	// Explicit routes (no nested "" paths)
	r.Post("/deals", d.createDeal)
	r.Get("/deals", d.listDeals)

	r.Post("/orders", d.createOrder)

	// SSE
	r.Get("/events", d.Hub.ServeHTTP)

	return r
}
