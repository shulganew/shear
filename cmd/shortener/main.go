package main

import (
	"net/http"

	"github.com/shulganew/shear.git/internal/config"
	webhandl "github.com/shulganew/shear.git/internal/web/handlers"
	"github.com/shulganew/shear.git/internal/web/router"
)

func main() {

	configApp := config.InitConfig()
	err := http.ListenAndServe(configApp.StartAddress, router.RouteShear(*webhandl.NewHandler(configApp)))
	if err != nil {
		panic(err)
	}

}
