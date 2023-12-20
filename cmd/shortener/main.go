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
)

func main() {

	app.InitLog()

	ctx, cancel := app.InitContext()
	defer cancel()

	conf := config.InitConfig()

	var db *sql.DB

	if conf.IsDB {
		db = storage.InitDB(ctx, conf.DSN)
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

	if err := http.ListenAndServe(conf.Address, router.RouteShear(conf, stor, db)); err != nil {
		panic(err)
	}

}
