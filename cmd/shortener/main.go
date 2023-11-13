package main

import (
	"net/http"

	"github.com/shulganew/shear.git/internal/app/initapp"
	"github.com/shulganew/shear.git/internal/web/router"
)

func main() {

	configApp := initapp.InitApp()

	err := http.ListenAndServe(configApp.StartAddress, router.RouteShear())
	if err != nil {
		panic(err)
	}

}
