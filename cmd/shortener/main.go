package main

import (
	"context"
	"database/sql"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/shulganew/shear.git/internal/app"
	"github.com/shulganew/shear.git/internal/config"
	"github.com/shulganew/shear.git/internal/service"
	"github.com/shulganew/shear.git/internal/storage"
	"github.com/shulganew/shear.git/internal/web/handlers"
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

	// сигнальный канал для завершения горутин
	doneCh := make(chan struct{})
	// закрываем его при завершении программы
	defer close(doneCh)

	// канал с данными
	inputCh := make(chan handlers.DelBrief, 1000)

	go func(ctx context.Context, s *service.StorageURL) {
		service := service.NewService(s)
		for res := range inputCh {
			service.DelelteBatch(ctx, res.UserID, res.Briefs)
		}

	}(ctx, stor)

	if err := http.ListenAndServe(conf.Address, router.RouteShear(conf, stor, db, inputCh)); err != nil {
		panic(err)
	}

}
