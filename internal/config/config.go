package config

import (
	"flag"
	"os"

	"github.com/shulganew/shear.git/internal/service"
	"github.com/shulganew/shear.git/internal/storage"
	"github.com/shulganew/shear.git/internal/web/validators"
	"go.uber.org/zap"
)

const DefaultHost string = "localhost:8080"

type Shear struct {
	//flag -a
	Address string
	//env var, or flag -b if env not exist
	Response string

	Backup service.Backup

	Storage storage.StorageURL
}

func InitConfig() *Shear {

	config := Shear{}

	//set logger
	InitLog()

	//read command line argue

	startAddress := flag.String("a", "localhost:8080", "start server address and port")
	resultAddress := flag.String("b", "localhost:8080", "answer address and port")
	tempf := flag.String("f", "", "Location of dump file")
	flag.Parse()
	//check and parse URL

	startaddr, startport := validators.CheckURL(*startAddress)
	answaddr, answport := validators.CheckURL(*resultAddress)

	//save config
	config.Address = startaddr + ":" + startport
	config.Response = answaddr + ":" + answport
	zap.S().Infoln("Server address: ", config.Address)

	//read OS ENV
	env, exist := os.LookupEnv(("SERVER_ADDRESS"))

	//if env var does not exist  - set def value
	if exist {
		config.Response = env
		zap.S().Infoln("Set result address from evn SERVER_ADDRESS: ", config.Response)

	} else {
		zap.S().Infoln("Env var SERVER_ADDRESS not found, use default", config.Response)
	}

	//define backup file
	config.Backup = service.Backup{}

	temp, exist := os.LookupEnv(("FILE_STORAGE_PATH"))
	if exist {
		config.Backup = *service.New(temp, true)
		zap.S().Infoln("Found backup's evn, use file: ", config.Backup.File)
	} else if *tempf != "" {
		config.Backup = *service.New(*tempf, true)
	} else {
		config.Backup = *service.New(*tempf, false)
	}

	zap.S().Infoln("Backup isActive: ", config.Backup.IsActive)

	//set MemoryStorage storage
	config.Storage = &storage.MemoryStorage{StoreURLs: []storage.Short{}}

	//load all dump links
	shorts, err := config.Backup.Load()
	if err != nil {
		zap.S().Error("Error load backup!", err)
	}

	//set MemoryStorage storage
	config.Storage = &storage.MemoryStorage{StoreURLs: shorts}

	zap.S().Infoln("Configuration complite")
	return &config
}

func (c *Shear) SetConfig(startAddress, resultAddress string) {
	c.Address = startAddress
	c.Response = resultAddress
}

func (c *Shear) GetStartAddr() string {
	return c.Address
}

func (c *Shear) GetResultAddr() string {
	return c.Response
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
