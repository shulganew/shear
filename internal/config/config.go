package config

import (
	"context"
	"database/sql"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/shulganew/shear.git/internal/service"
	"github.com/shulganew/shear.git/internal/storage"
	"github.com/shulganew/shear.git/internal/web/validators"
	"go.uber.org/zap"
)

const DefaultHost string = "localhost:8080"

type App struct {
	//flag -a
	Address string
	//env var, or flag -b if env not exist
	Response string

	Backup service.Backup

	Storage storage.StorageURL

	//data base connection
	DB *sql.DB
}

func InitConfig() (*App, context.CancelFunc, *sql.DB) {

	config := App{}

	//init Context
	ctx, cancel := InitContext()

	//set logger
	InitLog()

	//read command line argue

	startAddress := flag.String("a", "localhost:8080", "start server address and port")
	resultAddress := flag.String("b", "localhost:8080", "answer address and port")
	tempf := flag.String("f", "", "Location of dump file")
	dsnf := flag.String("d", "", "Data Source Name for DataBase connection")
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

	//init Storage
	dsn, exist := os.LookupEnv(("DATABASE_DSN"))
	//init shotrage DB from env
	if exist {
		zap.S().Infoln("Use DataBase Storge. Found location DATABASE_DSN, use: ", dsn)
		//init shotrage DB from env
		db := InitDB(ctx, dsn)
		if err := db.Ping(); err != nil {
			zap.S().Errorln("Error connect to DB from env: ", err)
		}
		//set sqk.DB to config
		config.DB = db

		//set Storage to DB
		config.Storage = &storage.DB{DB: db}

	} else if *dsnf != "" {
		dsn = *dsnf
		zap.S().Infoln("Use DataBase Storge. Found -d flag, use: ", dsn)
		//init shotrage DB from flag
		db := InitDB(ctx, dsn)
		if err := db.Ping(); err != nil {
			zap.S().Errorln("Error connect to DB from flag: ", err)
		}
		//set sqk.DB to config
		config.DB = db

		//set Storage to DB
		config.Storage = &storage.DB{DB: db}

	} else {
		//init memory storage
		zap.S().Infoln("Use Memory Storge.")
		config.Storage = &storage.Memory{}
	}

	//define backup file

	temp, exist := os.LookupEnv(("FILE_STORAGE_PATH"))
	if exist {
		config.Backup = *service.NewBackup(temp, true)
		zap.S().Infoln("Found backup's evn, use file: ", config.Backup.File)
	} else if *tempf != "" {
		config.Backup = *service.NewBackup(*tempf, true)
	} else {
		config.Backup = *service.NewBackup(*tempf, false)
	}

	zap.S().Infoln("Backup isActive: ", config.Backup.IsActive)

	//load all dump links
	shorts, err := config.Backup.Load()
	if err != nil {
		zap.S().Error("Error load backup!", err)
	}

	//upload shorts to Storage
	config.Storage.SetAll(ctx, shorts)

	//activate backup

	if config.Backup.IsActive {
		InitBackup(ctx, &config)
	} else {
		//if back is not active, make exit after ctrl+C
		go func() {
			<-ctx.Done()
			os.Exit(0)
		}()
	}

	zap.S().Infoln("Configuration complite")
	return &config, cancel, config.DB
}

func (c *App) SetConfig(startAddress, resultAddress string) {
	c.Address = startAddress
	c.Response = resultAddress
}

func (c *App) GetStartAddr() string {
	return c.Address
}

func (c *App) GetResultAddr() string {
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
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	go func() {
		<-exit
		cancel()

	}()
	return
}

// Activate backup
func InitBackup(ctx context.Context, config *App) {
	//Time machine
	service.TimeBackup(ctx, config.Storage, config.Backup)
	//backup on graceful
	service.Shutdown(ctx, config.Storage, config.Backup)
}

// Init Database
func InitDB(ctx context.Context, dsn string) (db *sql.DB) {

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		panic(err)
	}

	//create table short if not exist

	_, err = db.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS short (id SERIAL , brief TEXT NOT NULL, origin TEXT NOT NULL UNIQUE)")
	if err != nil {
		panic(err)
	}

	return
}

func InitMemoryStorge() *storage.StorageURL {
	return nil
}
