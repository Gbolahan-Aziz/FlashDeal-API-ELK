package store

import (
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"time"
	"github.com/google/uuid"	
	_ "github.com/mattn/go-sqlite3"

	"FlashDeal-API-ELK/internal/domain"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrInsufficient = errors.New("insufficient stock")
)

type SQLStore struct {
	db *sql.DB
}

// NewSQL opens (or creates) an SQLite db and runs migrations.
func NewSQL(dsn string) (*SQLStore, error) {
	// Ensure parent directory exists: e.g. /data for /data/flashdeals.db
	dir := filepath.Dir(dsn)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	if err := migrate(db); err != nil {
		return nil, err
	}
	return &SQLStore{db: db}, nil
}

func migrate(db *sql.DB) error {
	schema := `
CREATE TABLE IF NOT EXISTS deals (
  id        TEXT PRIMARY KEY,
  title     TEXT NOT NULL,
  price     REAL NOT NULL,
  stock     INTEGER NOT NULL,
  active    INTEGER NOT NULL,
  created_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS orders (
  id         TEXT PRIMARY KEY,
  deal_id    TEXT NOT NULL,
  qty        INTEGER NOT NULL,
  status     TEXT NOT NULL,
  created_at TEXT NOT NULL,
  FOREIGN KEY (deal_id) REFERENCES deals(id)
);
`
	_, err := db.Exec(schema)
	return err
}

// CreateDeal implements Store.CreateDeal
func (s *SQLStore) CreateDeal(nd domain.NewDeal) (*domain.Deal, error) {
	id := uuid.New().String() 
	now := time.Now().UTC()

	deal := &domain.Deal{
		ID:        id,
		Title:     nd.Title,
		Price:     nd.Price,
		Stock:     nd.Stock,
		Active:    nd.Active,
		CreatedAt: now,
	}

	_, err := s.db.Exec(
		`INSERT INTO deals (id, title, price, stock, active, created_at)
         VALUES (?, ?, ?, ?, ?, ?)`,
		deal.ID, deal.Title, deal.Price, deal.Stock, boolToInt(deal.Active), deal.CreatedAt.Format(time.RFC3339Nano),
	)
	if err != nil {
		return nil, err
	}
	return deal, nil
}

// ListDeals implements Store.ListDeals
func (s *SQLStore) ListDeals() ([]domain.Deal, error) {
	rows, err := s.db.Query(
		`SELECT id, title, price, stock, active, created_at
         FROM deals ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.Deal
	for rows.Next() {
		var d domain.Deal
		var activeInt int
		var created string
		if err := rows.Scan(&d.ID, &d.Title, &d.Price, &d.Stock, &activeInt, &created); err != nil {
			return nil, err
		}
		d.Active = activeInt == 1
		t, _ := time.Parse(time.RFC3339Nano, created)
		d.CreatedAt = t
		out = append(out, d)
	}
	return out, rows.Err()
}

// Order implements Store.Order
func (s *SQLStore) Order(dealID string, qty int) (*domain.Order, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var stock int
	err = tx.QueryRow(
		`SELECT stock FROM deals WHERE id = ? AND active = 1`,
		dealID,
	).Scan(&stock)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if stock < qty {
		return nil, ErrInsufficient
	}

	_, err = tx.Exec(
		`UPDATE deals SET stock = stock - ? WHERE id = ?`,
		qty, dealID,
	)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	order := &domain.Order{
		ID:      uuid.New().String(),
		DealID:  dealID,
		Qty:     qty,
		Status:  "accepted",
		Created: now,
	}

	_, err = tx.Exec(
		`INSERT INTO orders (id, deal_id, qty, status, created_at)
         VALUES (?, ?, ?, ?, ?)`,
		order.ID, order.DealID, order.Qty, order.Status, order.Created.Format(time.RFC3339Nano),
	)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return order, nil
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
