package main

import (
	"net/http"

	"github.com/shulganew/shear.git/internal/config"
	"github.com/shulganew/shear.git/internal/web/router"
)

func main() {

	configApp := config.InitConfig()
	err := http.ListenAndServe(configApp.StartAddress, router.RouteShear(configApp))
	if err != nil {
		panic(err)
	}

}
