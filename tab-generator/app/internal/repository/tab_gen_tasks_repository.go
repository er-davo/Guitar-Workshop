package repository

import (
	"context"
	"time"

	"tabgen/internal/models"

	sq "github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TabGenTaskRepository struct {
	db     *pgxpool.Pool
	getter *trmpgx.CtxGetter
	psql   sq.StatementBuilderType
}

func NewTabGenTaskRepository(db *pgxpool.Pool, c *trmpgx.CtxGetter) *TabGenTaskRepository {
	return &TabGenTaskRepository{
		db:     db,
		getter: c,
		psql:   sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *TabGenTaskRepository) Create(ctx context.Context, task *models.TabGenTask) error {
	if task.CreatedAt.IsZero() {
		task.CreatedAt = time.Now()
	}

	query := r.psql.Insert("tab_gen_tasks").
		Columns(
			"audio_sep_task_id",
			"status",
			"input_audio_name",
			"result_tab_name",
			"saved",
			"error",
			"created_at",
		).
		Values(
			task.AudioSepTaskID,
			string(task.Status),
			task.InputAudioName,
			task.ResultTabName,
			task.Saved,
			task.Error,
			task.CreatedAt,
		).Suffix("RETURNING id")

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	conn := r.getter.DefaultTrOrDB(ctx, r.db)
	err = conn.QueryRow(ctx, sql, args...).Scan(&task.ID)
	return wrapDBError(err)
}

func (r *TabGenTaskRepository) Get(ctx context.Context, id string) (*models.TabGenTask, error) {
	query := r.psql.Select(
		"id",
		"audio_sep_task_id",
		"status",
		"input_audio_name",
		"result_tab_name",
		"saved",
		"error",
		"created_at",
	).From("tab_gen_tasks").
		Where(sq.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	conn := r.getter.DefaultTrOrDB(ctx, r.db)
	task := &models.TabGenTask{}
	err = conn.QueryRow(ctx, sql, args...).Scan(
		&task.ID,
		&task.AudioSepTaskID,
		&task.Status,
		&task.InputAudioName,
		&task.ResultTabName,
		&task.Saved,
		&task.Error,
		&task.CreatedAt,
	)

	return task, wrapDBError(err)
}

func (r *TabGenTaskRepository) GetByAudioSepTaskID(ctx context.Context, id string) (*models.TabGenTask, error) {
	query := r.psql.Select(
		"id",
		"audio_sep_task_id",
		"status",
		"input_audio_name",
		"result_tab_name",
		"saved",
		"error",
		"created_at",
	).From("tab_gen_tasks").
		Where(sq.Eq{"audio_sep_task_id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	conn := r.getter.DefaultTrOrDB(ctx, r.db)
	task := &models.TabGenTask{}
	err = conn.QueryRow(ctx, sql, args...).Scan(
		&task.ID,
		&task.AudioSepTaskID,
		&task.Status,
		&task.InputAudioName,
		&task.ResultTabName,
		&task.Saved,
		&task.Error,
		&task.CreatedAt,
	)

	return task, wrapDBError(err)
}

func (r *TabGenTaskRepository) MarkSaved(ctx context.Context, id string) error {
	query := r.psql.Update("tab_gen_tasks").
		Set("saved", true).
		Where(sq.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return wrapDBError(err)
	}

	conn := r.getter.DefaultTrOrDB(ctx, r.db)
	_, err = conn.Exec(ctx, sql, args...)
	return wrapDBError(err)
}

func (r *TabGenTaskRepository) UpdateStatus(
	ctx context.Context,
	id string,
	status models.Status,
	errMsg *string,
	resultPath *string,
) error {
	query := r.psql.
		Update("tab_gen_tasks").
		Set("status", string(status)).
		Where(sq.Eq{"id": id})

	if errMsg != nil {
		query = query.Set("error", errMsg)
	}
	if resultPath != nil {
		query = query.Set("result_tab_path", resultPath)
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return wrapDBError(err)
	}

	conn := r.getter.DefaultTrOrDB(ctx, r.db)
	_, err = conn.Exec(ctx, sql, args...)
	return wrapDBError(err)
}

func (r *TabGenTaskRepository) TryUpdateStatus(
	ctx context.Context,
	id string,
	statusFrom []models.Status,
	to models.Status,
	errMsg *string,
	resultTabName *string,
) (bool, error) {
	if len(statusFrom) == 0 {
		return false, nil
	}

	fromStatuses := make([]string, 0, len(statusFrom))
	for _, s := range statusFrom {
		fromStatuses = append(fromStatuses, string(s))
	}

	query := r.psql.
		Update("tab_gen_tasks").
		Set("status", string(to)).
		Where(sq.Eq{
			"id": id,
		}).
		Where(sq.Eq{
			"status": fromStatuses,
		})

	if errMsg != nil {
		query = query.Set("error", errMsg)
	}

	if resultTabName != nil {
		query = query.Set("result_tab_name", resultTabName)
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return false, wrapDBError(err)
	}

	conn := r.getter.DefaultTrOrDB(ctx, r.db)
	cmdTag, err := conn.Exec(ctx, sql, args...)
	if err != nil {
		return false, wrapDBError(err)
	}

	if cmdTag.RowsAffected() == 0 {
		return false, nil
	}

	return true, nil
}
