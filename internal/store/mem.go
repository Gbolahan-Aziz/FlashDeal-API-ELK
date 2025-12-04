package store

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"FlashDeal-API-ELK/internal/domain"
)


type Mem struct {
	mu    sync.RWMutex
	deals map[string]*domain.Deal
}

func NewMemStore() *Mem {
	return &Mem{deals: make(map[string]*domain.Deal)}
}

func (m *Mem) CreateDeal(in domain.NewDeal) (*domain.Deal, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	d := &domain.Deal{
		ID: uuid.New().String(), Title: in.Title, Price: in.Price, Stock: in.Stock,
		Active: in.Active, CreatedAt: time.Now(),
	}
	m.deals[d.ID] = d
	return d, nil
}

func (m *Mem) ListDeals() ([]domain.Deal, error) {
	m.mu.RLock(); defer m.mu.RUnlock()
	out := make([]domain.Deal, 0, len(m.deals))
	for _, d := range m.deals {
		out = append(out, *d)
	}
	return out, nil
}

func (m *Mem) Order(dealID string, qty int) (*domain.Order, error) {
	m.mu.Lock(); defer m.mu.Unlock()
	d, ok := m.deals[dealID]
	if !ok || !d.Active {
		return nil, ErrNotFound
	}
	if qty <= 0 || d.Stock < qty {
		return nil, ErrInsufficient
	}
	d.Stock -= qty
	o := &domain.Order{
		ID: uuid.New().String(), DealID: d.ID, Qty: qty, Status: "accepted", Created: time.Now(),
	}
	return o, nil
}
