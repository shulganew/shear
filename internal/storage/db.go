package storage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/shulganew/shear.git/internal/service"
	"go.uber.org/zap"
)

type DB struct {
	master *sql.DB
}

func NewDB(ctx context.Context, master *sql.DB) (*DB, error) {
	db := DB{master: master}
	err := db.Start(ctx)
	return &db, err
}

func (base *DB) Set(ctx context.Context, userID, brief, origin string) error {

	err := base.master.QueryRowContext(ctx, "INSERT INTO short (user_id, brief, origin) VALUES ($1, $2, $3) ", userID, brief, origin).Scan()
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

func (base *DB) GetOrigin(ctx context.Context, brief string) (origin string, ok bool) {
	row := base.master.QueryRowContext(ctx, "SELECT origin from short where brief=$1", brief)
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

func (base *DB) GetBrief(ctx context.Context, origin string) (brief string, ok bool) {
	row := base.master.QueryRowContext(ctx, "SELECT brief from short where origin=$1", origin)
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

func (base *DB) GetAll(ctx context.Context) []service.Short {

	rows, err := base.master.QueryContext(ctx, "SELECT id, user_id, brief, origin from short")
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	shorts := []service.Short{}
	for rows.Next() {
		var short service.Short
		err = rows.Scan(&short.ID, &short.UUID, &short.Brief, &short.Origin)
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

func (base *DB) GetUserAll(ctx context.Context, userID string) []service.Short {

	rows, err := base.master.QueryContext(ctx, "SELECT id, user_id, brief, origin FROM short WHERE user_id=$1", userID)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	shorts := []service.Short{}
	for rows.Next() {
		var short service.Short
		err = rows.Scan(&short.ID, &short.UUID, &short.Brief, &short.Origin)
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

func (base *DB) SetAll(ctx context.Context, shorts []service.Short) error {

	tx, err := base.master.Begin()
	if err != nil {
		panic(err)
	}

	prep, err := base.master.PrepareContext(ctx, "INSERT INTO short (user_id, brief, origin) VALUES ($1, $2, $3)")
	if err != nil {
		panic(err)
	}

	for _, short := range shorts {
		zap.S().Infoln("Insert!!!!!!!!!!!!!! ", short.Brief)
		_, err := prep.ExecContext(ctx, short.UUID, short.Brief, short.Origin)
		if err != nil {
			var pgErr *pgconn.PgError
			// if URL exist in DataBase
			if errors.As(err, &pgErr) && pgerrcode.UniqueViolation == pgErr.Code {
				//get brief string
				if brief, ok := base.GetBrief(ctx, short.Origin); ok {

					return NewErrDuplicatedShort(short.SessionID, brief, short.Origin, pgErr)
				}
			}
			tx.Rollback()
			return err
		}

	}
	err = tx.Commit()
	if err != nil {
		panic(err)
	}
	return nil
}

func (base *DB) DelelteBatch(ctx context.Context, userID string, briefs []string) {
	zap.S().Infoln("UPDATE !!!!!!!!!!!!!! ", briefs)
	time.Sleep(time.Second)
	tx, err := base.master.Begin()

	if err != nil {
		panic(err)
	}

	//prerare bulck request to database
	//fille
	userIDs := make([]string, len(briefs), len(briefs))
	for i := range briefs {
		userIDs[i] = userID
	}

	bulck := `
	UPDATE short SET is_deleted = TRUE 
	FROM (SELECT unnest($1::text[]) AS user_id, unnest($2::text[]) AS brief) AS data_table 
	WHERE short.user_id = data_table.user_id AND short.brief = data_table.brief;
	`

	_, err = base.master.ExecContext(ctx, bulck, userIDs, briefs)

	if err != nil {
		tx.Rollback()
		panic(err)
	}

	err = tx.Commit()
	if err != nil {
		panic(err)
	}

}

// Init Database
func InitDB(ctx context.Context, dsn string) (db *sql.DB, err error) {

	db, err = sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	//create table short if not exist

	_, err = db.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS short (id SERIAL , user_id TEXT NULL, brief TEXT NOT NULL, origin TEXT NOT NULL UNIQUE)")
	if err != nil {
		return nil, err
	}

	//upgrade table if uuid not exist
	_, err = db.ExecContext(ctx, "ALTER TABLE short ADD COLUMN IF NOT EXISTS user_id TEXT")
	if err != nil {
		return nil, err
	}

	//upgrade table if is_deleted not exist
	_, err = db.ExecContext(ctx, "ALTER TABLE short ADD COLUMN IF NOT EXISTS is_deleted BOOLEAN DEFAULT FALSE")
	if err != nil {
		return nil, err
	}

	return
}

func (base *DB) Start(ctx context.Context) error {
	// ждем 3 секунды - если не смогли стартовать - возвращаем ошибку
	// читаем контекст
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	err := base.master.PingContext(ctx)
	defer cancel()
	return err
}
