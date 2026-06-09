//go:build !cgo

package store

import (
	"fmt"
	"FlashDeal-API-ELK/internal/domain"
)

// SQLStore is a placeholder type for non-CGO builds.
type SQLStore struct{}

// NewSQL is unavailable without CGO. Callers should fall back to in-memory.
func NewSQL(_ string) (*SQLStore, error) {
	return nil, fmt.Errorf("SQLite requires CGO_ENABLED=1; falling back to in-memory store")
}

func (s *SQLStore) CreateDeal(_ domain.NewDeal) (*domain.Deal, error) { return nil, nil }
func (s *SQLStore) ListDeals() ([]domain.Deal, error)                 { return nil, nil }
func (s *SQLStore) Order(_ string, _ int) (*domain.Order, error)      { return nil, nil }
