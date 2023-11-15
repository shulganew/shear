package main

import (
	"net/http"

	"github.com/shulganew/shear.git/internal/app"
	webhandl "github.com/shulganew/shear.git/internal/web/handlers"
	"github.com/shulganew/shear.git/internal/web/router"
)

func main() {

	configApp := app.Init()
	handler := webhandl.URLHandler{}

	handler.SetStorage(configApp.Storage)
	handler.SetConfig(configApp)
	err := http.ListenAndServe(configApp.StartAddress, router.RouteShear(handler))
	if err != nil {
		panic(err)
	}

}
