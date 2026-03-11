package repository

import (
	"context"

	"tab-service/internal/models"

	sq "github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TabRepository struct {
	db     *pgxpool.Pool
	getter *trmpgx.CtxGetter
	psql   sq.StatementBuilderType
}

func NewTabRepository(db *pgxpool.Pool, c *trmpgx.CtxGetter) *TabRepository {
	return &TabRepository{
		db:     db,
		getter: c,
		psql:   sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *TabRepository) Create(ctx context.Context, tab *models.Tab) error {
	if tab == nil {
		return ErrNilValue
	}

	query, args, err := r.psql.
		Insert("tabs").
		Columns("name", "file_path").
		Values(tab.Name, tab.Path).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return err
	}

	conn := r.getter.DefaultTrOrDB(ctx, r.db)
	err = conn.QueryRow(ctx, query, args...).Scan(&tab.ID)
	return wrapDBError(err)
}

func (r *TabRepository) Delete(ctx context.Context, id string) error {
	if id == "" {
		return ErrInvalidID
	}

	query, args, err := r.psql.
		Delete("tabs").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return err
	}

	conn := r.getter.DefaultTrOrDB(ctx, r.db)
	cmd, err := conn.Exec(ctx, query, args...)
	if err != nil {
		return wrapDBError(err)
	}
	if cmd.RowsAffected() == 0 {
		return ErrNoRowsAffected
	}
	return nil
}

func (r *TabRepository) MarkDeleted(ctx context.Context, id string) error {
	if id == "" {
		return ErrInvalidID
	}

	query, args, err := r.psql.
		Update("tabs").
		Set("deleted_at", sq.Expr("NOW()")).
		Where(sq.And{
			sq.Eq{"id": id},
			sq.Expr("deleted_at IS NULL"),
		}).
		ToSql()
	if err != nil {
		return err
	}

	conn := r.getter.DefaultTrOrDB(ctx, r.db)
	cmd, err := conn.Exec(ctx, query, args...)
	if err != nil {
		return wrapDBError(err)
	}
	if cmd.RowsAffected() == 0 {
		return ErrNoRowsAffected
	}

	return nil
}

func (r *TabRepository) Get(ctx context.Context, id string) (*models.Tab, error) {
	if id == "" {
		return nil, ErrInvalidID
	}

	query, args, err := r.psql.
		Select(
			"id", "name", "file_path",
			"created_at", "deleted_at",
		).
		From("tabs").
		Where(sq.And{
			sq.Eq{"id": id},
			sq.Expr("deleted_at IS NULL"),
		}).
		ToSql()
	if err != nil {
		return nil, err
	}

	tab := &models.Tab{}
	conn := r.getter.DefaultTrOrDB(ctx, r.db)
	err = conn.QueryRow(ctx, query, args...).Scan(
		&tab.ID, &tab.Name, &tab.Path,
		&tab.CreatedAt, &tab.DeletedAt,
	)
	if err != nil {
		return nil, wrapDBError(err)
	}

	return tab, nil
}

func (r *TabRepository) Search(
	ctx context.Context,
	name string,
	limit uint64,
	offset uint64,
) ([]*models.Tab, error) {

	query, args, err := r.psql.
		Select("id", "name", "file_path", "created_at").
		Prefix("WITH q AS (SELECT websearch_to_tsquery('simple', ?) AS query)", name).
		From("tabs, q").
		Where("deleted_at IS NULL").
		Where("search_vector @@ q.query").
		OrderBy("ts_rank_cd(search_vector, q.query) DESC, created_at DESC").
		Limit(limit).
		Offset(offset).
		ToSql()
	if err != nil {
		return nil, err
	}

	conn := r.getter.DefaultTrOrDB(ctx, r.db)

	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return nil, wrapDBError(err)
	}
	defer rows.Close()

	var tabs []*models.Tab

	for rows.Next() {
		tab := &models.Tab{}
		if err := rows.Scan(
			&tab.ID,
			&tab.Name,
			&tab.Path,
			&tab.CreatedAt,
		); err != nil {
			return nil, wrapDBError(err)
		}
		tabs = append(tabs, tab)
	}

	return tabs, rows.Err()
}
