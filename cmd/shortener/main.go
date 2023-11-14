package main

import (
	"net/http"

	"github.com/shulganew/shear.git/internal/app"
	"github.com/shulganew/shear.git/internal/storage"
	webhandl "github.com/shulganew/shear.git/internal/web/handlers"
	"github.com/shulganew/shear.git/internal/web/router"
)

func main() {

	configApp := app.Init()
	handler := webhandl.URLHandler{}
	handler.SetMapStorage(storage.NewMapStorage())
	handler.SetConfig(configApp)
	err := http.ListenAndServe(configApp.StartAddress, router.RouteShear(handler))
	if err != nil {
		panic(err)
	}

}
