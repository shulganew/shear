package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

type DB struct {
	DB *sql.DB
}

func (base *DB) Set(ctx context.Context, brief, origin string) error {

	err := base.DB.QueryRowContext(ctx, "INSERT INTO short (brief, origin) VALUES ($1, $2) ", brief, origin).Scan()
	if err != nil {
		zap.S().Infoln("Insert error!: ", origin)
		var pgErr *pgconn.PgError
		// if URL exist in DataBase
		if errors.As(err, &pgErr) && pgerrcode.UniqueViolation == pgErr.Code {
			//get brief string
			if brief, ok := base.GetBrief(ctx, origin); ok {
				zap.S().Infoln("Found duplicated URL: ", origin)
				return NewErrDuplicatedURL(brief, origin, pgErr)
			}
		}

		// if URL exist in DataBase
		if err == sql.ErrNoRows {
			//insert - no rows returned
			return nil
		}
		return err
	}
	return nil
}

func (base *DB) GetBrief(ctx context.Context, origin string) (brief string, ok bool) {
	row := base.DB.QueryRowContext(ctx, "SELECT brief from short where origin=$1", origin)
	err := row.Scan(&brief)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", false
		}
		panic(err)
	}
	err = row.Err()
	if err != nil {
		panic(err)
	}
	return brief, true
}

func (base *DB) GetOrigin(ctx context.Context, brief string) (origin string, ok bool) {
	row := base.DB.QueryRowContext(ctx, "SELECT origin from short where brief=$1", brief)
	err := row.Scan(&origin)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", false
		}
		panic(err)
	}
	err = row.Err()
	if err != nil {
		panic(err)
	}
	return origin, true
}

func (base *DB) GetAll(ctx context.Context) []Short {

	rows, err := base.DB.QueryContext(ctx, "SELECT id, brief, origin from short")
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	shorts := []Short{}
	for rows.Next() {
		var short Short
		err = rows.Scan(&short.ID, &short.Brief, &short.Origin)
		if err != nil {
			panic(err)
		}

		shorts = append(shorts, short)
	}

	err = rows.Err()
	if err != nil {
		panic(err)
	}

	return shorts

}

func (base *DB) SetAll(ctx context.Context, shorts []Short) error {

	tx, err := base.DB.Begin()
	if err != nil {
		panic(err)
	}

	prep, err := base.DB.PrepareContext(ctx, "INSERT INTO short (brief, origin) VALUES ($1, $2)")
	if err != nil {
		panic(err)
	}

	for _, short := range shorts {
		_, err := prep.ExecContext(ctx, short.Brief, short.Origin)
		if err != nil {
			var pgErr *pgconn.PgError
			// if URL exist in DataBase
			if errors.As(err, &pgErr) && pgerrcode.UniqueViolation == pgErr.Code {
				//get brief string
				if brief, ok := base.GetBrief(ctx, short.Origin); ok {

					return NewErrDuplicatedURL(brief, short.Origin, pgErr)
				}
			}
			return err
		}

	}
	err = tx.Commit()
	if err != nil {
		panic(err)
	}
	return nil
}
