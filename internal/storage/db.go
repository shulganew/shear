package storage

import (
	"context"
	"database/sql"
	"strings"

	"go.uber.org/zap"
)

type DB struct {
	DB *sql.DB
}

func (base *DB) Set(ctx context.Context, brief, origin string) Short {
	var lastInsertID int
	err := base.DB.QueryRowContext(ctx, "INSERT INTO short (brief, origin) VALUES ($1, $2) RETURNING id", brief, origin).Scan(&lastInsertID)
	if err != nil {
		zap.S().Errorln("Error inserting new short in database.")
		panic(err)
	}
	return Short{ID: lastInsertID, Brief: brief, Origin: origin}
}

func (base *DB) GetBrief(ctx context.Context, origin string) (brief string, ok bool) {
	row := base.DB.QueryRowContext(ctx, "SELECT brief from short where brief=$1", brief)
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
	return
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
	return
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

func (base *DB) SetAll(ctx context.Context, shorts []Short) {
	for _, short := range shorts {
		_, err := base.DB.ExecContext(ctx, "INSERT INTO short (id, brief, origin) VALUES ($1, $2, $3) RETURNING id", short.ID, short.Brief, short.Origin)
		if err != nil {
			//duplicate key value, alredy exist
			if strings.Contains(err.Error(), "SQLSTATE 23505") {
				continue
			}
			panic(err)
		}

	}
}
