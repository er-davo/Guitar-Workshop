package repository

import (
	"context"

	"api-gateway/internal/models"

	"github.com/jackc/pgx/v5"
)

type TabRepository interface {
	Create(ctx context.Context, tab *models.Tab) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*models.Tab, error)
	FindByNameLike(ctx context.Context, name string) ([]*models.Tab, error)
}

type tabRepository struct {
	db *pgx.Conn
}

func NewTabRepository(db *pgx.Conn) *tabRepository {
	return &tabRepository{db: db}
}

func (r *tabRepository) Create(ctx context.Context, tab *models.Tab) error {
	query := `INSERT INTO tabs (name, file_path) VALUES ($1, $2) RETURNING id`
	return r.db.QueryRow(ctx, query, pgx.QueryExecModeSimpleProtocol, tab.Name, tab.Path).Scan(&tab.ID)
}

func (r *tabRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM tabs WHERE id = $1`
	_, err := r.db.Exec(ctx, query, pgx.QueryExecModeSimpleProtocol, id)
	return err
}

func (r *tabRepository) GetByID(ctx context.Context, id string) (*models.Tab, error) {
	tab := new(models.Tab)
	query := `SELECT id, name, file_path FROM tabs WHERE id = $1`
	err := r.db.QueryRow(ctx, query, pgx.QueryExecModeSimpleProtocol, id).Scan(&tab.ID, &tab.Name, &tab.Path)
	if err != nil {
		return nil, err
	}
	return tab, nil
}

func (r *tabRepository) FindByNameLike(ctx context.Context, name string) ([]*models.Tab, error) {
	query := `SELECT id, name, file_path FROM tabs WHERE name ILIKE '%' || $1 || '%'`
	rows, err := r.db.Query(ctx, query, pgx.QueryExecModeSimpleProtocol, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tabs []*models.Tab
	for rows.Next() {
		tab := new(models.Tab)
		if err := rows.Scan(&tab.ID, &tab.Name, &tab.Path); err != nil {
			return nil, err
		}
		tabs = append(tabs, tab)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tabs, nil
}
