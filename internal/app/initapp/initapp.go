package initapp

import (
	"flag"
	"log"
	"os"

	"github.com/shulganew/shear.git/internal/app/config"
	utils "github.com/shulganew/shear.git/internal/core"
)

// Init app parameters from cmd and env
func InitApp() *config.ConfigShear {
	configApp := config.GetConfig()

	//init config
	//read command line argue

	startAddress := flag.String("a", "localhost:8080", "start server address and port")
	resultAddress := flag.String("b", "localhost:8080", "answer address and port")

	flag.Parse()
	//check and parse URL
	startaddr, startport := utils.CheckAddress(*startAddress)
	answaddr, answport := utils.CheckAddress(*resultAddress)

	//save config
	configApp.StartAddress = startaddr + ":" + startport
	configApp.ResultAddress = answaddr + ":" + answport
	log.Println("Server address: ", configApp.StartAddress)

	//read OS ENV
	envAddress, exist := os.LookupEnv(("SERVER_ADDRESS"))

	//if env var does not exist - set def value
	if exist {
		configApp.ResultAddress = envAddress
		log.Println("Set result address from evn SERVER_ADDRESS: ", configApp.ResultAddress)

	} else {
		log.Println("Env var SERVER_ADDRESS not found, use default", configApp.ResultAddress)
	}

	log.Println("Config main: ", configApp)
	return configApp
}
