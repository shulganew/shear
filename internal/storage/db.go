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

	err := base.master.QueryRowContext(ctx, "INSERT INTO short (user_id, brief, origin, is_deleted) VALUES ($1, $2, $3, $4) ", userID, brief, origin, false).Scan()
	if err != nil {
		//zap.S().Infoln("Insert error!: ", origin)

		var pgErr *pgconn.PgError
		// if URL exist in DataBase
		if errors.As(err, &pgErr) && pgerrcode.UniqueViolation == pgErr.Code {

			//get brief string
			if brief, ok, _ := base.GetBrief(ctx, origin); ok {

				//check if marked as deleted - recreate!
				if base.isDeleted(ctx, brief) {
					base.Recover(ctx, userID, brief)
					return nil
				}

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

func (base *DB) GetOrigin(ctx context.Context, brief string) (origin string, existed bool, isDeleted bool) {

	row := base.master.QueryRowContext(ctx, "SELECT id, user_id, brief, origin, is_deleted FROM short WHERE brief=$1", brief)

	var short service.Short
	err := row.Scan(&short.ID, &short.UUID, &short.Brief, &short.Origin, &short.IsDeleted)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", false, false
		}
		panic(err)
	}

	return short.Origin, true, short.IsDeleted
}

func (base *DB) GetBrief(ctx context.Context, origin string) (brief string, existed bool, isDeleted bool) {
	row := base.master.QueryRowContext(ctx, "SELECT id, user_id, brief, origin, is_deleted FROM short WHERE origin=$1", origin)

	var short service.Short
	err := row.Scan(&short.ID, &short.UUID, &short.Brief, &short.Origin, &short.IsDeleted)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", false, false
		}
		panic(err)
	}

	return short.Brief, true, short.IsDeleted

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
		_, err := prep.ExecContext(ctx, short.UUID, short.Brief, short.Origin)
		if err != nil {
			var pgErr *pgconn.PgError
			// if URL exist in DataBase
			if errors.As(err, &pgErr) && pgerrcode.UniqueViolation == pgErr.Code {
				//get brief string
				if brief, ok, _ := base.GetBrief(ctx, short.Origin); ok {

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
	zap.S().Infoln("Delete briefs: ", len(briefs))
	//prerare bulck request to database
	//fille
	userIDs := make([]string, len(briefs))
	for i := range briefs {
		userIDs[i] = userID
	}

	bulck := `
	UPDATE short SET is_deleted = TRUE 
	FROM (SELECT unnest($1::text[]) AS user_id, unnest($2::text[]) AS brief) AS data_table 
	WHERE short.user_id = data_table.user_id AND short.brief = data_table.brief;
	`

	_, err := base.master.ExecContext(ctx, bulck, userIDs, briefs)

	if err != nil {
		panic(err)
	}

}

func (base *DB) isDeleted(ctx context.Context, brief string) bool {

	row := base.master.QueryRowContext(ctx, "SELECT is_deleted FROM short WHERE brief=$1", brief)

	var isDeleted bool
	err := row.Scan(&isDeleted)

	if err != nil {
		panic(err)
	}
	zap.S().Infoln("isDeleted: ", brief, isDeleted)
	return isDeleted

}

func (base *DB) Recover(ctx context.Context, userID string, brief string) {

	_, err := base.master.ExecContext(ctx, "UPDATE short SET is_deleted=FALSE, user_id=$1 WHERE brief=$2", userID, brief)
	zap.S().Infoln("Recover!!!", userID, brief)
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
