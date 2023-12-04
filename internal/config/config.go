package config

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

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

func InitConfig(ctx context.Context) *Shear {

	config := Shear{}

	//set logger
	InitLog()

	//set MemoryStorage storage
	config.Storage = &storage.MemoryStorage{StoreURLs: []storage.Short{}}

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
		config.Backup = *service.New(temp, true, config.Storage)
		zap.S().Infoln("Found backup's evn, use file: ", config.Backup.File)
	} else if *tempf != "" {
		config.Backup = *service.New(*tempf, true, config.Storage)
	} else {
		config.Backup = *service.New(*tempf, false, config.Storage)
	}

	zap.S().Infoln("Backup isActive: ", config.Backup.IsActive)

	//load all dump links
	shorts, err := config.Backup.Load()
	if err != nil {
		zap.S().Error("Error load backup!", err)
	}

	//set MemoryStorage storage
	config.Storage = &storage.MemoryStorage{StoreURLs: shorts}

	//activate backup
	if config.Backup.IsActive {

		//Time machine
		service.TimeBackup(config.Storage, config.Backup)
		//backup on graceful
		service.Shutdown(ctx, config.Storage, config.Backup)
	}

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

// Init context from graceful shutdown. Send to all function for return by syscall.SIGINT, syscall.SIGTERM
func InitContext() (ctx context.Context, cancel context.CancelFunc) {
	exit := make(chan os.Signal, 1)
	ctx, cancel = context.WithCancel(context.Background())
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-exit
		cancel()
	}()
	fmt.Println("End Init")
	return
}
