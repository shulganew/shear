// Package app use for initializing main functionalities of application Shortener:
//
// • initialize DB Postgres or memory storage;
//
// • initialize backup if enable;
//
// • graceful shutdown realization;
//
// • initialize zap logger;
package app

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/shulganew/shear.git/internal/config"
	"github.com/shulganew/shear.git/internal/service"
	"github.com/shulganew/shear.git/internal/storage"
	"go.uber.org/zap"
)

// Function InitApp initialize Database, Backup and create Delete service.
func InitApp(ctx context.Context, conf config.Config, db *sql.DB, delCh chan service.DelBatch) (*service.Shorten, *service.Backup, *service.Delete) {
	var stor service.StorageURL
	var err error
	// load storage
	if conf.IsDB() && db != nil {
		// use db storage
		stor, err = storage.NewDB(ctx, db)
		if err != nil {
			zap.S().Errorln("Error connect to DB from env: ", err)
			// use memory storage
			stor = storage.NewMemory()
			zap.S().Infoln("Use memory storage: database not pinging")
		}
		zap.S().Infoln("Use database storage")
	} else {
		// use memory storage
		stor = storage.NewMemory()
		zap.S().Infoln("Use memory storage")
	}

	short := service.NewService(stor)

	var backup *service.Backup
	// define backup file
	if conf.IsBackup() {
		backup = InitBackup(ctx, short, conf.GetBackupPath())
		zap.S().Infoln("Backup activated, path: ", conf.GetBackupPath())

		// load all dump links
		shorts, err := backup.Load()
		if err != nil {
			zap.S().Error("Error load backup!", err)
		}

		// upload shorts to Storage
		stor.SetAll(ctx, shorts)
	}

	del := service.NewDelete(delCh, &conf)
	zap.S().Infoln("Application init complete")
	return short, backup, del
}

// Init context from graceful shutdown. Send to all function for return by
//
//	syscall.SIGINT, syscall.SIGTERM
func InitContext() (ctx context.Context, cancel context.CancelFunc) {
	ctx, cancel = signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	return
}

// Print app version data on startup.
func Intro(version string, date string, commit string) {
	fmt.Println("Build version: ", version)
	fmt.Println("Build date: ", date)
	fmt.Println("Build commit: ", commit)
}

// Initialization of a zap logger.
func InitLog() zap.SugaredLogger {
	logger, err := zap.NewDevelopment()
	if err != nil {

		panic(err)
	}
	zap.ReplaceGlobals(logger)
	defer logger.Sync()
	sugar := *logger.Sugar()
	defer sugar.Sync()
	return sugar
}

// Init Database connection using pgx driver.
func InitDB(ctx context.Context, dsn string) (db *sql.DB, err error) {
	db, err = sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	return
}

// Activate backup
func InitBackup(ctx context.Context, short *service.Shorten, file string) *service.Backup {
	backup := &service.Backup{File: file}
	// Time machine.
	service.TimeBackup(ctx, short, *backup)
	return backup
}
