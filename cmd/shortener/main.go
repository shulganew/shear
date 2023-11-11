package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/shulganew/shear.git/internal/app/config"
	webhandl "github.com/shulganew/shear.git/internal/handlers"
)

func RouteShear() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/{id}", webhandl.GetUrl)
	r.Post("/", webhandl.SetUrl)

	return r
}

// Init app parameters from cmd and env
func initApp() *config.ConfigShear {
	configApp := config.GetConfig()
	//init config
	//read command line argue

	startAddress := flag.String("a", "localhost:8080", "start address and port")
	resultAddress := flag.String("b", "localhost:8080", "start address and port")

	flag.Parse()
	configApp.StartAddress = *startAddress
	log.Println("Server address: ", configApp.StartAddress)
	//read OS ENV
	envAddress, exist := os.LookupEnv(("SERVER_ADDRESS"))

	//if env var does not exist - set def value
	if exist {
		log.Println("Set result address from evn SERVER_ADDRESS: ", resultAddress)
		configApp.ResultAddress = envAddress
	} else {
		log.Println("Env var SERVER_ADDRESS not found, use default localhost:8080.")
		configApp.ResultAddress = config.DefaultHost
	}

	log.Println("Config main: ", configApp.StartAddress)
	return configApp
}

func main() {

	configApp := initApp()
	err := http.ListenAndServe(configApp.StartAddress, RouteShear())
	if err != nil {
		panic(err)
	}

}
