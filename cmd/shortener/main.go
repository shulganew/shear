package main

import (
	"net/http"

	"github.com/shulganew/shear.git/internal/app/initapp"
	"github.com/shulganew/shear.git/internal/storage"
	webhandl "github.com/shulganew/shear.git/internal/web/handlers"
	"github.com/shulganew/shear.git/internal/web/router"
)

func main() {

	configApp := initapp.InitApp()
	handler := webhandl.URLHandler{}
	handler.SetMapStorage(&storage.MapStorage{})
	err := http.ListenAndServe(configApp.StartAddress, router.RouteShear(handler))
	if err != nil {
		panic(err)
	}

}
