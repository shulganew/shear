package storage

import (
	"context"
	"database/sql/driver"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDBGetAll(t *testing.T) {
	tests := []struct {
		name string
		rows [][]driver.Value
	}{
		{
			name: "Base SQL test: GetAll",
			rows: [][]driver.Value{{0, uuid.NewString(), "dzafbfsx", "http://yandex1.ru/"}, {1, uuid.NewString(), "dzafbfsy", "http://yandex2.ru/"}},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			sql, sqlmock, err := sqlmock.New()
			if err != nil {
				t.Errorf("An error '%s' was not expected when opening a stub database connection", err)
			}

			// columns are prefixed with "o" since we used sqlstruct to generate them
			columns := []string{"o_id", "o_user_id", "o_brief", "o_origin"}
			sqlmock.ExpectQuery("SELECT id, user_id, brief, origin from short").
				//WithArgs(0).
				//WillReturnRows(sqlmock.NewRows(columns).FromCSVString("95, 018dea9b-7085-75f5-91c5-2ba674052348, QoFwCenk, http://ya144016.localhost"))
				WillReturnRows(sqlmock.NewRows(columns).AddRows(tt.rows[0], tt.rows[1]))
			sqlmock.ExpectClose()

			db, err := NewDB(ctx, sql)
			assert.NoError(t, err)
			shorts := db.GetAll(ctx)
			for i, short := range shorts {
				assert.Equal(t, tt.rows[i][0], short.ID)
				assert.Equal(t, tt.rows[i][1], short.UserID.String)
				assert.Equal(t, tt.rows[i][2], short.Brief)
				assert.Equal(t, tt.rows[i][3], short.Origin)
			}

			// db.Close() ensures that all expectations have been met
			if err = sql.Close(); err != nil {
				t.Errorf("Error '%s' was not expected while closing the database", err)
			}
		})

	}
}

func TestDBGetUserAll(t *testing.T) {
	tests := []struct {
		name   string
		userID string
		rows   [][]driver.Value
	}{
		{
			name:   "Base SQL test: GetUserAll",
			userID: "018dea9b-7085-75f5-91c5-2ba674052348",
			rows:   [][]driver.Value{{0, "018dea9b-7085-75f5-91c5-2ba674052348", "dzafbfsx", "http://yandex1.ru/"}, {1, "018dea9b-7085-75f5-91c5-2ba674052348", "dzafbfsy", "http://yandex2.ru/"}},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			sql, sqlmock, err := sqlmock.New()
			if err != nil {
				t.Errorf("An error '%s' was not expected when opening a stub database connection", err)
			}

			// columns are prefixed with "o" since we used sqlstruct to generate them
			columns := []string{"o_id", "o_user_id", "o_brief", "o_origin"}
			sqlmock.ExpectQuery("SELECT id, user_id, brief, origin FROM short WHERE user_id=(.+)").
				WillReturnRows(sqlmock.NewRows(columns).AddRows(tt.rows[0], tt.rows[1]))
			sqlmock.ExpectClose()

			db, err := NewDB(ctx, sql)
			assert.NoError(t, err)
			shorts := db.GetUserAll(ctx, tt.userID)
			for i, short := range shorts {
				assert.Equal(t, short.ID, tt.rows[i][0])
				assert.Equal(t, short.UserID.String, tt.rows[i][1])
				assert.Equal(t, short.Brief, tt.rows[i][2])
				assert.Equal(t, short.Origin, tt.rows[i][3])
			}

			// db.Close() ensures that all expectations have been met
			if err = sql.Close(); err != nil {
				t.Errorf("Error '%s' was not expected while closing the database", err)
			}
		})

	}
}

func TestDBGetOrigin(t *testing.T) {
	tests := []struct {
		name    string
		brief   string
		origin  string
		row     []driver.Value
		norows  bool
		existed bool
	}{
		{
			name:    "Base SQL test: GetOrigin1",
			brief:   "dfgjhdfj",
			origin:  "http://yandex.ru",
			row:     []driver.Value{0, uuid.NewString(), "dzafbfsx", "http://yandex.ru", false},
			existed: true,
			norows:  false,
		},
		{
			name:    "Base SQL test: GetOrigin2",
			brief:   "dfgjhdfj",
			origin:  "",
			row:     []driver.Value{0, uuid.NewString(), "dzafbfsx", "http://yandex.ru", false},
			existed: false,
			norows:  true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			sql, sqlmock, err := sqlmock.New()
			if err != nil {
				t.Errorf("An error '%s' was not expected when opening a stub database connection", err)
			}

			// columns are prefixed with "o" since we used sqlstruct to generate them
			columns := []string{"o_id", "o_user_id", "o_brief", "o_origin", "o_is_deleted"}

			if !tt.norows {
				sqlmock.ExpectQuery("SELECT id, user_id, brief, origin, is_deleted FROM short WHERE brief=(.+)").
					WillReturnRows(sqlmock.NewRows(columns).AddRows(tt.row))
			} else {
				sqlmock.ExpectQuery("SELECT id, user_id, brief, origin, is_deleted FROM short WHERE brief=(.+)").
					WillReturnRows(sqlmock.NewRows(columns))
			}
			sqlmock.ExpectClose()

			db, err := NewDB(ctx, sql)
			assert.NoError(t, err)
			origin, exist, isDel := db.GetOrigin(ctx, tt.brief)

			assert.Equal(t, tt.origin, origin)
			assert.Equal(t, tt.existed, exist)
			assert.False(t, isDel)

			// db.Close() ensures that all expectations have been met
			if err = sql.Close(); err != nil {
				t.Errorf("Error '%s' was not expected while closing the database", err)
			}
		})

	}
}

func TestDBGetBrief(t *testing.T) {
	tests := []struct {
		name    string
		brief   string
		origin  string
		row     []driver.Value
		norows  bool
		existed bool
	}{
		{
			name:    "Base SQL test: Get Brief 1",
			brief:   "dzafbfsx",
			origin:  "http://yandex.ru",
			row:     []driver.Value{0, uuid.NewString(), "dzafbfsx", "http://yandex.ru", false},
			existed: true,
			norows:  false,
		},
		{
			name:    "Base SQL test: GetBrief 2",
			brief:   "",
			origin:  "http://yandex.ru",
			row:     []driver.Value{0, uuid.NewString(), "dzafbfsx", "http://yandex.ru", false},
			existed: false,
			norows:  true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			sql, sqlmock, err := sqlmock.New()
			if err != nil {
				t.Errorf("An error '%s' was not expected when opening a stub database connection", err)
			}

			// columns are prefixed with "o" since we used sqlstruct to generate them
			columns := []string{"o_id", "o_user_id", "o_brief", "o_origin", "o_is_deleted"}

			if !tt.norows {
				sqlmock.ExpectQuery("SELECT id, user_id, brief, origin, is_deleted FROM short WHERE origin(.+)").
					WillReturnRows(sqlmock.NewRows(columns).AddRows(tt.row))
			} else {
				sqlmock.ExpectQuery("SELECT id, user_id, brief, origin, is_deleted FROM short WHERE origin(.+)").
					WillReturnRows(sqlmock.NewRows(columns))
			}
			sqlmock.ExpectClose()

			db, err := NewDB(ctx, sql)
			assert.NoError(t, err)
			brief, exist, isDel := db.GetBrief(ctx, tt.origin)

			assert.Equal(t, tt.brief, brief)
			assert.Equal(t, tt.existed, exist)
			assert.False(t, isDel)

			// db.Close() ensures that all expectations have been met
			if err = sql.Close(); err != nil {
				t.Errorf("Error '%s' was not expected while closing the database", err)
			}
		})

	}
}

func TestDBSet(t *testing.T) {
	tests := []struct {
		name    string
		userID  string
		brief   string
		origin  string
		row     []driver.Value
		norows  bool
		existed bool
	}{
		{
			name:    "Base SQL test: Set 1",
			userID:  "018dea9b-7085-75f5-91c5-2ba674052348",
			brief:   "dzafbfsx",
			origin:  "http://yandex.ru",
			row:     []driver.Value{0, uuid.NewString(), "dzafbfsx", "http://yandex.ru", false},
			existed: false,
			norows:  false,
		},
		{
			name:    "Base SQL test: Set with error 2",
			userID:  "018dea9b-7085-75f5-91c5-2ba674052348",
			brief:   "dzafbfsx",
			origin:  "http://yandex.ru",
			row:     []driver.Value{0, uuid.NewString(), "dzafbfsx", "http://yandex.ru", false},
			existed: true,
			norows:  false,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			sql, sqlmock, err := sqlmock.New()
			if err != nil {
				t.Errorf("An error '%s' was not expected when opening a stub database connection", err)
			}

			if !tt.existed {
				sqlmock.ExpectQuery("INSERT INTO short \\(user_id, brief, origin, is_deleted\\) VALUES (.+)").
					WillReturnRows(sqlmock.NewRows([]string{}))
			} else {
				sqlmock.ExpectQuery("INSERT INTO short \\(user_id, brief, origin, is_deleted\\) VALUES (.+)").
					WillReturnError(NewErrDuplicatedURL(tt.brief, tt.origin, errors.New("Duplicated.")))
			}
			sqlmock.ExpectClose()

			db, err := NewDB(ctx, sql)
			assert.NoError(t, err)
			err = db.Set(ctx, tt.userID, tt.brief, tt.origin)

			// If rows duplicated.
			if !tt.existed {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}

			// db.Close() ensures that all expectations have been met
			if err = sql.Close(); err != nil {
				t.Errorf("Error '%s' was not expected while closing the database", err)
			}
		})

	}
}

func TestDBIsDeleted(t *testing.T) {
	tests := []struct {
		row     driver.Value
		name    string
		brief   string
		existed bool
	}{
		{
			name:    "Base SQL test: Is Deleted false",
			brief:   "dfgjhdfj",
			row:     true,
			existed: true,
		},

		{
			name:    "Base SQL test: Is Deleted true",
			brief:   "dfgjhdfj",
			row:     false,
			existed: false,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			sql, sqlmock, err := sqlmock.New()
			if err != nil {
				t.Errorf("An error '%s' was not expected when opening a stub database connection", err)
			}

			// columns are prefixed with "o" since we used sqlstruct to generate them
			columns := []string{"o_is_deleted"}

			if !tt.existed {
				sqlmock.ExpectQuery("SELECT is_deleted FROM short WHERE brief=(.+)").
					WillReturnRows(sqlmock.NewRows(columns).AddRow(tt.row))
			} else {
				sqlmock.ExpectQuery("SELECT is_deleted FROM short WHERE brief=(.+)").
					WillReturnRows(sqlmock.NewRows(columns).AddRow(tt.row))
			}

			sqlmock.ExpectClose()

			db, err := NewDB(ctx, sql)
			assert.NoError(t, err)
			isDel := db.isDeleted(ctx, tt.brief)

			if !tt.existed {
				assert.False(t, isDel)
			} else {
				assert.True(t, isDel)
			}

			// db.Close() ensures that all expectations have been met
			if err = sql.Close(); err != nil {
				t.Errorf("Error '%s' was not expected while closing the database", err)
			}
		})
	}
}

func TestDRRecover(t *testing.T) {
	tests := []struct {
		name    string
		userID  string
		brief   string
		existed bool
	}{
		{
			name:    "Base SQL test: Is Deleted false",
			userID:  "018dea9b-7085-75f5-91c5-2ba674052348",
			brief:   "dfgjhdfj",
			existed: false,
		},
		{
			name:    "Base SQL test: Is Deleted false",
			userID:  "018dea9b-7085-75f5-91c5-2ba674052348",
			brief:   "dfgjhdfj",
			existed: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			sql, smock, err := sqlmock.New()
			if err != nil {
				t.Errorf("An error '%s' was not expected when opening a stub database connection", err)
			}
			if !tt.existed {
				smock.ExpectExec("(.+)").WithArgs(tt.userID, tt.brief).
					WillReturnResult(sqlmock.NewResult(0, 1)) // no insert id, 1 affected row
			} else {
				smock.ExpectExec("(.+)").WithArgs(tt.userID, tt.brief).
					WillReturnError(errors.New("Not updated"))
			}

			smock.ExpectClose()

			db, err := NewDB(ctx, sql)
			assert.NoError(t, err)

			err = db.Recover(ctx, tt.userID, tt.brief)
			if !tt.existed {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
			// db.Close() ensures that all expectations have been met
			if err = sql.Close(); err != nil {
				t.Errorf("Error '%s' was not expected while closing the database", err)
			}
		})
	}
}
