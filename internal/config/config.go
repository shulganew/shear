package config

import (
	"flag"
	"net/url"
	"os"

	"github.com/shulganew/shear.git/internal/storage"
	"github.com/shulganew/shear.git/internal/web/validators"
	"go.uber.org/zap"
)

const DefaultHost string = "localhost:8080"

type Shear struct {
	//flag -a
	StartAddress string
	//env var, or flag -b if env not exist
	ResultAddress string

	Storage storage.StorageURL
}

func InitConfig() *Shear {

	config := Shear{}

	//set logger
	InitLog()

	//read command line argue

	startAddress := flag.String("a", "localhost:8080", "start server address and port")
	resultAddress := flag.String("b", "localhost:8080", "answer address and port")
	flag.Parse()
	//check and parse URL

	startaddr, startport := validators.CheckURL(*startAddress)
	answaddr, answport := validators.CheckURL(*resultAddress)

	//save config
	config.StartAddress = startaddr + ":" + startport
	config.ResultAddress = answaddr + ":" + answport
	zap.S().Infoln("Server address: ", config.StartAddress)

	//read OS ENV
	envAddress, exist := os.LookupEnv(("SERVER_ADDRESS"))

	//if env var does not exist - set def value
	if exist {
		config.ResultAddress = envAddress
		zap.S().Infoln("Set result address from evn SERVER_ADDRESS: ", config.ResultAddress)

	} else {
		zap.S().Infoln("Env var SERVER_ADDRESS not found, use default", config.ResultAddress)
	}

	//set Map storage
	config.Storage = &storage.MapStorage{StoreURLs: make(map[string]url.URL)}

	return &config
}

func (c *Shear) SetConfig(startAddress, resultAddress string) {
	c.StartAddress = startAddress
	c.ResultAddress = resultAddress
}

func (c *Shear) GetStartAddr() string {
	return c.StartAddress
}

func (c *Shear) GetResultAddr() string {
	return c.ResultAddress
}

func InitLog() zap.SugaredLogger {
	logger, err := zap.NewDevelopment()
	if err != nil {

		panic(err)
	}
	zap.ReplaceGlobals(logger)
	defer logger.Sync()

	sugar := *logger.Sugar()

	defer sugar.Sync()
	return sugar
}
