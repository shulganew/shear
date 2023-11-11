package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/shulganew/shear.git/internal/app/config"
	webhandl "github.com/shulganew/shear.git/internal/handlers"
)

func Router() chi.Router {
	r := chi.NewRouter()
	r.Get("/{id}", webhandl.GetUrl)
	r.Post("/", webhandl.SetUrl)

	return r
}

func main() {

	//read command line argue
	startAddress := flag.String("a", "localhost:8080", "start address and port")
	flag.Parse()

	//read OS ENV
	resultAddress, exist := os.LookupEnv(("SERVER_ADDRESS"))

	//if env var does not exist - set def value
	if exist {
		log.Println("Set result address from evn SERVER_ADDRESS: ", resultAddress)
	} else {
		resultAddress = config.DefaultHost
		log.Println("Env var SERVER_ADDRESS not found, use default.")
	}

	configApp := config.GetConfig()
	//init config
	configApp.StartAddress = *startAddress
	configApp.ResultAddress = resultAddress

	log.Println("Config main: ", configApp.StartAddress)

	err := http.ListenAndServe(configApp.StartAddress, Router())
	if err != nil {
		panic(err)
	}

}
