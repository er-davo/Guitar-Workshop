package repository

import (
	"context"

	"api-gateway/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TabRepository interface {
	Create(ctx context.Context, tab *models.Tab) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*models.Tab, error)
	FindByNameLike(ctx context.Context, name string) ([]*models.Tab, error)
}

type tabRepository struct {
	db *pgxpool.Pool
}

func NewTabRepository(db *pgxpool.Pool) *tabRepository {
	return &tabRepository{db: db}
}

func (r *tabRepository) Create(ctx context.Context, tab *models.Tab) error {
	if tab == nil {
		return ErrNilValue
	}
	query := `INSERT INTO tabs (name, file_path) VALUES ($1, $2) RETURNING id`
	err := r.db.QueryRow(
		ctx,
		query,
		pgx.QueryExecModeSimpleProtocol,
		tab.Name,
		tab.Path,
	).Scan(&tab.ID)
	return wrapDBError(err)
}

func (r *tabRepository) Delete(ctx context.Context, id string) error {
	if id == "" {
		return ErrInvalidID
	}
	query := `DELETE FROM tabs WHERE id = $1`
	cmd, err := r.db.Exec(ctx, query, pgx.QueryExecModeSimpleProtocol, id)
	if cmd.RowsAffected() == 0 {
		return ErrNoRowsAffected
	}
	return wrapDBError(err)
}

func (r *tabRepository) GetByID(ctx context.Context, id string) (*models.Tab, error) {
	if id == "" {
		return nil, ErrInvalidID
	}
	tab := new(models.Tab)
	query := `SELECT id, name, file_path FROM tabs WHERE id = $1`
	err := r.db.QueryRow(ctx, query, pgx.QueryExecModeSimpleProtocol, id).Scan(&tab.ID, &tab.Name, &tab.Path)
	if err != nil {
		return nil, wrapDBError(err)
	}
	return tab, nil
}

func (r *tabRepository) FindByNameLike(ctx context.Context, name string) ([]*models.Tab, error) {
	query := `SELECT id, name, file_path FROM tabs WHERE name ILIKE '%' || $1 || '%'`
	rows, err := r.db.Query(ctx, query, pgx.QueryExecModeSimpleProtocol, name)
	if err != nil {
		return nil, wrapDBError(err)
	}
	defer rows.Close()

	var tabs []*models.Tab
	for rows.Next() {
		tab := new(models.Tab)
		if err := rows.Scan(&tab.ID, &tab.Name, &tab.Path); err != nil {
			return nil, wrapDBError(err)
		}
		tabs = append(tabs, tab)
	}

	if err := rows.Err(); err != nil {
		return nil, wrapDBError(err)
	}

	return tabs, nil
}
