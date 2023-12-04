package main

import (
	"net/http"

	"github.com/shulganew/shear.git/internal/config"
	"github.com/shulganew/shear.git/internal/web/router"
)

func main() {

	configApp := config.InitConfig()
	//activate backup
	if configApp.Backup.IsActive {
		ctx, cancel := config.InitContext()
		config.InitBackup(ctx, configApp)
		defer cancel()

	}

	err := http.ListenAndServe(configApp.Address, router.RouteShear(configApp))
	if err != nil {
		panic(err)
	}

}
