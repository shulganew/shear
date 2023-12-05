package main

import (
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/shulganew/shear.git/internal/config"
	"github.com/shulganew/shear.git/internal/web/router"
)

func main() {

	configApp, cancel, db := config.InitConfig()
	//defer close context
	if cancel != nil {
		defer cancel()
	}

	//defer close db
	defer db.Close()

	if err := http.ListenAndServe(configApp.Address, router.RouteShear(configApp)); err != nil {
		panic(err)
	}

}
