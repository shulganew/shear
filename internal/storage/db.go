// Package implements storage layer of a Shortener for database and memory storage.
package storage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/shulganew/shear.git/internal/entities"
	"go.uber.org/zap"
)

// Database with connection filed and set of Repo messages, implements shortener interface.
type DB struct {
	master *sql.DB
}

// Constructor of Database obj.
func NewDB(ctx context.Context, master *sql.DB) (*DB, error) {
	db := DB{master: master}
	err := db.Start(ctx)
	return &db, err
}

// Add short and original URL to storage.
func (base *DB) Add(ctx context.Context, userID, brief, origin string) error {
	err := base.master.QueryRowContext(ctx, "INSERT INTO short (user_id, brief, origin, is_deleted) VALUES ($1, $2, $3, $4) ", userID, brief, origin, false).Scan()
	if err != nil {
		var pgErr *pgconn.PgError
		// if URL exist in DataBase
		if errors.As(err, &pgErr) && pgerrcode.UniqueViolation == pgErr.Code {
			// get brief string
			if brief, ok, _ := base.GetBrief(ctx, origin); ok {

				// check if marked as deleted - recreate!
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
			// insert - no rows returned
			return nil
		}
		return err
	}
	return nil
}

// Set all user's short and original URLs from Short slice.
func (base *DB) AddAll(ctx context.Context, shorts []entities.Short) error {
	tx, err := base.master.Begin()
	if err != nil {
		zap.S().Errorln(err)
	}

	prep, err := base.master.PrepareContext(ctx, "INSERT INTO short (user_id, brief, origin) VALUES ($1, $2, $3)")
	if err != nil {
		zap.S().Errorln(err)
	}

	for _, short := range shorts {
		_, err = prep.ExecContext(ctx, short.UserID, short.Brief, short.Origin)
		if err != nil {
			var pgErr *pgconn.PgError
			// if URL exist in DataBase
			if errors.As(err, &pgErr) && pgerrcode.UniqueViolation == pgErr.Code {
				// get brief string
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
		zap.S().Errorln(err)
	}
	return nil
}

// Get original URL from storage.
func (base *DB) GetOrigin(ctx context.Context, brief string) (origin string, existed bool, isDeleted bool) {
	row := base.master.QueryRowContext(ctx, "SELECT id, user_id, brief, origin, is_deleted FROM short WHERE brief=$1", brief)
	var short entities.Short
	err := row.Scan(&short.ID, &short.UserID, &short.Brief, &short.Origin, &short.IsDeleted)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", false, false
		}
		zap.S().Errorln(err)
	}
	return short.Origin, true, short.IsDeleted
}

// Get short URL from storage.
func (base *DB) GetBrief(ctx context.Context, origin string) (brief string, existed bool, isDeleted bool) {
	row := base.master.QueryRowContext(ctx, "SELECT id, user_id, brief, origin, is_deleted FROM short WHERE origin=$1", origin)
	var short entities.Short
	err := row.Scan(&short.ID, &short.UserID, &short.Brief, &short.Origin, &short.IsDeleted)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", false, false
		}
		zap.S().Errorln(err)
	}

	return short.Brief, true, short.IsDeleted
}

// Get all short and original URLs in Short slice.
func (base *DB) GetAll(ctx context.Context) []entities.Short {
	rows, err := base.master.QueryContext(ctx, "SELECT id, user_id, brief, origin from short")
	if err != nil {
		zap.S().Errorln(err)
	}

	defer rows.Close()

	shorts := []entities.Short{}
	for rows.Next() {
		var short entities.Short
		err = rows.Scan(&short.ID, &short.UserID, &short.Brief, &short.Origin)
		if err != nil {
			zap.S().Errorln(err)
		}

		shorts = append(shorts, short)
	}

	err = rows.Err()
	if err != nil {
		zap.S().Errorln(err)
	}

	return shorts
}

// Get all user's short and original URLs in Short slice.
func (base *DB) GetUserAll(ctx context.Context, userID string) []entities.Short {
	rows, err := base.master.QueryContext(ctx, "SELECT id, user_id, brief, origin FROM short WHERE user_id=$1", userID)
	if err != nil {
		zap.S().Errorln(err)
	}
	defer rows.Close()

	shorts := []entities.Short{}
	for rows.Next() {
		var short entities.Short
		err = rows.Scan(&short.ID, &short.UserID, &short.Brief, &short.Origin)
		if err != nil {
			zap.S().Errorln(err)
		}

		shorts = append(shorts, short)
	}
	err = rows.Err()
	if err != nil {
		zap.S().Errorln(err)
	}
	return shorts
}

// Mark all user's URLs by short URL in briefs slice.
func (base *DB) DeleteBatch(ctx context.Context, userID string, briefs []string) (err error) {
	// Prepare bulk request to database.
	userIDs := make([]string, len(briefs))
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
		zap.S().Errorln(err)
	}
	return
}

// Check if URL mark as deleted by short URL.
func (base *DB) isDeleted(ctx context.Context, brief string) bool {
	row := base.master.QueryRowContext(ctx, "SELECT is_deleted FROM short WHERE brief=$1", brief)

	var isDeleted bool
	err := row.Scan(&isDeleted)
	if err != nil {
		zap.S().Errorln(err)
	}
	zap.S().Infoln("isDeleted: ", brief, isDeleted)
	return isDeleted
}

// Undelete user's URL.
func (base *DB) Recover(ctx context.Context, userID string, brief string) error {
	_, err := base.master.ExecContext(ctx, "UPDATE short SET is_deleted=FALSE, user_id=$1 WHERE brief=$2", userID, brief)
	if err != nil {
		zap.S().Errorln(err)
		return err
	}
	return nil
}

// Connection to database with ping checking timeout (3 sec).
func (base *DB) Start(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	err := base.master.PingContext(ctx)
	defer cancel()
	return err
}

// Get totoal number of shorts.
func (base *DB) GetNumShorts(ctx context.Context) (num int, err error) {
	row := base.master.QueryRowContext(ctx, "SELECT COUNT(DISTINCT user_id) FROM short WHERE user_id IS NOT NULL")

	err = row.Scan(&num)
	if err != nil {
		return 0, err
	}
	zap.S().Infoln("Stat num of users: ", num)
	return num, nil
}

// Get totoal number of users.
func (base *DB) GetNumUsers(ctx context.Context) (num int, err error) {
	row := base.master.QueryRowContext(ctx, "SELECT COUNT(origin) FROM short")

	err = row.Scan(&num)
	if err != nil {
		return 0, err
	}
	zap.S().Infoln("Stat num of users: ", num)
	return num, nil
}
