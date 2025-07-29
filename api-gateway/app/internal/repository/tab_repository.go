package repository

import (
	"context"

	"api-gateway/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TabRepository struct {
	db *pgxpool.Pool
}

func NewTabRepository(db *pgxpool.Pool) *TabRepository {
	return &TabRepository{db: db}
}

func (r *TabRepository) Create(ctx context.Context, tab *models.Tab) error {
	query := `INSERT INTO tabs (name, file_path) VALUES ($1, $2) RETURNING id`
	return r.db.QueryRow(ctx, query, tab.Name, tab.Path).Scan(&tab.ID)
}

func (r *TabRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM tabs WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *TabRepository) GetByID(ctx context.Context, id int64) (*models.Tab, error) {
	tab := new(models.Tab)
	query := `SELECT id, name, file_path FROM tabs WHERE id = $1`
	err := r.db.QueryRow(ctx, query, id).Scan(&tab.ID, &tab.Name, &tab.Path)
	if err != nil {
		return nil, err
	}
	return tab, nil
}
