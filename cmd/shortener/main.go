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

	config := config.GetConfig()
	//init config
	config.StartAddress = *startAddress
	config.ResultAddress = resultAddress

	log.Println("Config main: ", config.StartAddress)
	r := chi.NewRouter()
	r.Get("/{id}", webhandl.GetUrl)
	r.Post("/", webhandl.SetUrl)

	err := http.ListenAndServe(config.StartAddress, r)
	if err != nil {
		panic(err)
	}

}
