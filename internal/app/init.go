package app

import (
	"context"
	"database/sql"
	"os"
	"os/signal"
	"syscall"

	"github.com/shulganew/shear.git/internal/config"
	"github.com/shulganew/shear.git/internal/service"
	"github.com/shulganew/shear.git/internal/storage"
	"go.uber.org/zap"
)

func InitApp(ctx context.Context, conf config.Config, db *sql.DB) (*service.StorageURL, *service.Backup) {


	//Storage
	var stor service.StorageURL

	//load storage
	if conf.IsDB && db != nil {
		if err := db.Ping(); err != nil {
			zap.S().Errorln("Error connect to DB from env: ", err)
			//use memory storage
			//set MemoryStorage storage
			stor = storage.NewMemory()
			zap.S().Infoln("Use memory storage: database not pinging")
		}
		//use db storage
		stor = storage.NewDB(db)
		zap.S().Infoln("Use database storage")
	} else {
		//use memory storage
		stor = storage.NewMemory()
		zap.S().Infoln("Use memory storage")
	}

	var backup *service.Backup
	//define backup file

	if conf.IsBackup {
		backup = service.InitBackup(ctx, stor, conf.BackupPath)
		zap.S().Infoln("Backup activated, path: ", conf.BackupPath)

		//load all dump links
		shorts, err := backup.Load()
		if err != nil {
			zap.S().Error("Error load backup!", err)
		}

		//upload shorts to Storage
		stor.SetAll(ctx, shorts)

	}

	zap.S().Infoln("Application init complite")
	return &stor, backup

}

// Init context from graceful shutdown. Send to all function for return by syscall.SIGINT, syscall.SIGTERM
func InitContext() (ctx context.Context, cancel context.CancelFunc) {
	exit := make(chan os.Signal, 1)
	ctx, cancel = context.WithCancel(context.Background())
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	go func() {
		<-exit
		cancel()

	}()
	return
}

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
