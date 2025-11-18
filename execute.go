package boilerplate

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func Get[T any](db *sqlx.DB, query string, args ...any) (t T, err error) {
	err = db.Get(&t, query, args...)
	log.Trace().Err(err).Any("result", t).Str("_query", query).Any("args", args).Msg("GET")
	return
}

func Select[T any](db *sqlx.DB, query string, args ...any) (t T, err error) {
	err = db.Select(&t, query, args...)
	log.Trace().Err(err).Any("result", t).Str("_query", query).Any("args", args).Msg("SELECT")
	return
}

func Exec(db *sqlx.DB, query string, args ...any) (err error) {
	r, err := db.Exec(query, args...)
	log.Trace().Err(err).Any("result", r).Str("_query", query).Any("args", args).Msg("EXEC")
	return
}

func NamedExecReturning(db *sqlx.DB, dest any, query string, args ...any) error {
	rows, err := db.NamedQuery(query, args)
	log.Trace().Err(err).Any("result", rows).Str("_query", query).Any("args", args).Msg("NAMED_EXEC_RET")

	if err != nil {
		return err
	}
	defer func() { _ = rows.Close() }()

	if rows.Next() {
		return rows.StructScan(dest)
	} else {
		if err := rows.Err(); err != nil {
			return err
		}
		return sql.ErrNoRows
	}
}

func NamedExec(db *sqlx.DB, dest any, query string, args ...any) error {
	rows, err := db.NamedExec(query, args)
	log.Trace().Err(err).Any("result", rows).Str("_query", query).Any("args", args).Msg("NAMED_EXEC")
	return err
}
