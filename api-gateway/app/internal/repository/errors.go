package repository

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrInvalidID           = errors.New("invalid id")
	ErrNilValue            = errors.New("nil value")
	ErrNotFound            = newProxyErr(pgx.ErrNoRows, "not found")
	ErrDuplicate           = errors.New("duplicate")
	ErrForeignKeyViolation = errors.New("foreign key violation")
	ErrNoRowsAffected      = errors.New("no rows affected")
	ErrNotNullViolation    = errors.New("not null violation")
	ErrInvalidDateTime     = errors.New("invalid datetime")
	ErrUndefinedColumn     = errors.New("undefined column")
	ErrUndefinedTable      = errors.New("undefined table")
)

func wrapDBError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": // unique_violation
			return ErrDuplicate
		case "23503": // foreign_key_violation
			return ErrForeignKeyViolation
		case "23502": // not_null_violation
			return ErrNotNullViolation
		case "22007", "22008": // invalid_datetime_format, datetime_field_overflow
			return ErrInvalidDateTime
		case "42703": // undefined_column
			return ErrUndefinedColumn
		case "42P01": // undefined_table
			return ErrUndefinedTable
		default:
			return fmt.Errorf("postgres error [%s]: %w", pgErr.Code, err)
		}
	}

	return err
}

type proxyError struct {
	msg        string
	background error
}

func newProxyErr(background error, msg string) error {
	return &proxyError{msg: msg, background: background}
}

func (p *proxyError) Error() string { return p.msg + ": " + p.background.Error() }
func (p *proxyError) Unwrap() error { return p.background }
