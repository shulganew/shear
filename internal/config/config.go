package config

import (
	"flag"
	"log"
	"os"

	"github.com/shulganew/shear.git/internal/web/netaddr"
)

const DefaultHost string = "localhost:8080"

type ConfigShear struct {
	//flag -a
	StartAddress string
	//env var, or flag -b if env not exist
	ResultAddress string
}

func InitConfig() *ConfigShear {

	config := ConfigShear{}
	//read command line argue

	startAddress := flag.String("a", "localhost:8080", "start server address and port")
	resultAddress := flag.String("b", "localhost:8080", "answer address and port")
	flag.Parse()
	//check and parse URL

	startaddr, startport := netaddr.CheckAddress(*startAddress)
	answaddr, answport := netaddr.CheckAddress(*resultAddress)

	//save config
	config.StartAddress = startaddr + ":" + startport
	config.ResultAddress = answaddr + ":" + answport
	log.Println("Server address: ", config.StartAddress)

	//read OS ENV
	envAddress, exist := os.LookupEnv(("SERVER_ADDRESS"))

	//if env var does not exist - set def value
	if exist {
		config.ResultAddress = envAddress
		log.Println("Set result address from evn SERVER_ADDRESS: ", config.ResultAddress)

	} else {
		log.Println("Env var SERVER_ADDRESS not found, use default", config.ResultAddress)
	}

	return &config
}

func (c *ConfigShear) SetConfig(startAddress, resultAddress string) {
	c.StartAddress = startAddress
	c.ResultAddress = resultAddress
}

func (c *ConfigShear) GetStartAddr() string {
	return c.StartAddress
}

func (c *ConfigShear) GetResultAddr() string {
	return c.ResultAddress
}
