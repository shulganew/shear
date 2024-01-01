package main

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"sync"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/shulganew/shear.git/internal/app"
	"github.com/shulganew/shear.git/internal/config"
	"github.com/shulganew/shear.git/internal/service"
	"github.com/shulganew/shear.git/internal/storage"
	"github.com/shulganew/shear.git/internal/web/router"
	"go.uber.org/zap"
)

func main() {

	app.InitLog()

	ctx, cancel := app.InitContext()
	defer cancel()

	conf := config.InitConfig()

	var db *sql.DB
	var err error
	if conf.IsDB {
		db, err = storage.InitDB(ctx, conf.DSN)
		if err != nil {
			db = nil
			conf.IsDB = false
			zap.S().Errorln("Can't connect to Database!", err)
		}
		defer db.Close()

	}
	stor, backup := app.InitApp(ctx, *conf, db)

	//Use fanIn pattern for storing data from delete requests
	finalCh := make(chan service.DelBatch)
	defer close(finalCh)

	var waitDel sync.WaitGroup

	go func(ctx context.Context, stor *service.StorageURL, finalCh chan service.DelBatch, wg *sync.WaitGroup) {
		serviceURL := service.NewService(stor)
		for {
			select {
			case <-ctx.Done():
				zap.S().Infoln("Waiting of update delete...")
				wg.Wait()
				if conf.IsBackup {
					service.Shutdown(*stor, *backup)
				}
				os.Exit(0)
			case delBatch := <-finalCh:
				serviceURL.DelelteBatch(ctx, delBatch)
			}
		}
	}(ctx, stor, finalCh, &waitDel)

	if err := http.ListenAndServe(conf.Address, router.RouteShear(conf, stor, db, finalCh, &waitDel)); err != nil {
		panic(err)
	}

}
