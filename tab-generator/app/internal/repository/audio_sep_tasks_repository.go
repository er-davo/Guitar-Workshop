package repository

import (
	"context"
	"time"

	"tabgen/internal/models"

	sq "github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AudioSepTaskRepository struct {
	db     *pgxpool.Pool
	getter *trmpgx.CtxGetter
	psql   sq.StatementBuilderType
}

func NewAudioSepTaskRepository(db *pgxpool.Pool, c *trmpgx.CtxGetter) *AudioSepTaskRepository {
	return &AudioSepTaskRepository{
		db:     db,
		getter: c,
		psql:   sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *AudioSepTaskRepository) Create(ctx context.Context, task *models.AudioSepTask) error {
	if task.ID == "" {
		return ErrInvalidID
	}
	if task.CreatedAt.IsZero() {
		task.CreatedAt = time.Now()
	}

	query, args, err := r.psql.
		Insert("audio_sep_tasks").
		Columns(
			"id",
			"status",
			"input_audio_name",
			"separated_dir_name",
			"error",
			"created_at",
		).
		Values(
			task.ID,
			string(task.Status),
			task.InputAudioName,
			task.SeparatedDirName,
			task.Error,
			task.CreatedAt,
		).
		ToSql()
	if err != nil {
		return wrapDBError(err)
	}

	conn := r.getter.DefaultTrOrDB(ctx, r.db)
	_, err = conn.Exec(ctx, query, args...)
	return wrapDBError(err)
}

func (r *AudioSepTaskRepository) Get(ctx context.Context, id string) (*models.AudioSepTask, error) {
	query, args, err := r.psql.
		Select(
			"id",
			"status",
			"input_audio_name",
			"separated_dir_name",
			"error",
			"created_at",
		).
		From("audio_sep_tasks").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}

	conn := r.getter.DefaultTrOrDB(ctx, r.db)
	task := &models.AudioSepTask{}
	err = conn.QueryRow(ctx, query, args...).Scan(
		&task.ID,
		&task.Status,
		&task.InputAudioName,
		&task.SeparatedDirName,
		&task.Error,
		&task.CreatedAt,
	)
	if err != nil {
		return nil, wrapDBError(err)
	}

	return task, nil
}
