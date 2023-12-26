package main

import (
	"database/sql"
	"net/http"
	"os"

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

	go func() {

		<-ctx.Done()
		if conf.IsBackup {
			service.Shutdown(*stor, *backup)
		}
		os.Exit(0)

	}()

	// // сигнальный канал для завершения горутин
	// doneCh := make(chan struct{})
	// // закрываем его при завершении программы
	// defer close(doneCh)

	// //var wg sync.WaitGroup
	// var mu = new(sync.Mutex)
	// var cond = sync.NewCond(mu)
	// chGen := concurrent.NewChgen(cond)

	// go func(ctx context.Context, dCh chan struct{}, s *service.StorageURL, cond *sync.Cond) {
	// 	mu.Lock()
	// 	defer mu.Unlock()
	// 	for {
	// 		cond.Wait()
	// 		finalCh := concurrent.FanIn(dCh, chGen.GetAllChannels())
	// 		for res := range finalCh {
	// 			zap.S().Infoln("Number of briefs in main: ", len(res.Briefs))
	// 			service := service.NewService(s)
	// 			service.DelelteBatch(ctx, res.UserID, res.Briefs)

	// 		}

	// 	}

	// }(ctx, doneCh, stor, cond)

	if err := http.ListenAndServe(conf.Address, router.RouteShear(conf, stor, db)); err != nil {
		panic(err)
	}

}
