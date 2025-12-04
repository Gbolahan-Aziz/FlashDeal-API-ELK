package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/rs/zerolog"
	"FlashDeal-API-ELK/internal/domain"
	"FlashDeal-API-ELK/internal/sse"
	"FlashDeal-API-ELK/internal/store"
)

type Store interface {
	CreateDeal(domain.NewDeal) (*domain.Deal, error)
	ListDeals() ([]domain.Deal, error)
	Order(dealID string, qty int) (*domain.Order, error)
}

type Deps struct {
	Log    zerolog.Logger
	Store  Store
	Hub    *sse.Hub
	Pub    interface{ Publish(context.Context, string, []byte) error }
	Config interface{ } // reserved
}

func (d Deps) createDeal(w http.ResponseWriter, r *http.Request) {
	var in domain.NewDeal
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest); return
	}
	deal, err := d.Store.CreateDeal(in)
	if err != nil { http.Error(w, "cannot create", 500); return }

	// broadcast
	rid := requestID(r)
	d.Hub.Broadcast(sse.JSONEvent("deal.created", rid, deal, time.Now().UnixMilli()))

	writeJSON(w, http.StatusCreated, deal)
}

func (d Deps) listDeals(w http.ResponseWriter, r *http.Request) {
	ds, _ := d.Store.ListDeals()
	writeJSON(w, http.StatusOK, ds)
}

func (d Deps) createOrder(w http.ResponseWriter, r *http.Request) {
	var in domain.NewOrder
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest); return
	}
	o, err := d.Store.Order(in.DealID, in.Qty)
	switch {
	case err == nil:
		rid := requestID(r)
		d.Hub.Broadcast(sse.JSONEvent("order.created", rid, o, time.Now().UnixMilli()))
		writeJSON(w, http.StatusCreated, o)
	case errors.Is(err, store.ErrNotFound):
		http.Error(w, "deal not found", http.StatusNotFound)
	case errors.Is(err, store.ErrInsufficient):
		http.Error(w, "insufficient stock", http.StatusConflict)
	default:
		http.Error(w, "cannot order", 500)
	}
}

func requestID(r *http.Request) string {
	if v := r.Header.Get("X-Request-ID"); v != "" { return v }
	return "" // middleware will set response header anyway
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
