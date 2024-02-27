// Package responsible for reading config from flag or ENV variables. After reading environment during init application config var locate in Config struct.
package config

import (
	"flag"
	"os"

	"github.com/shulganew/shear.git/internal/web/validators"
	"go.uber.org/zap"
)

// Default host.
const DefaultHost string = "localhost:8080"

// Struct for store main app config.
type Config struct {
	Address    string // flag -a
	Response   string //env var, or flag -b if env not exist
	IsBackup   bool   // is backup enable
	BackupPath string // file location Path for backup
	IsDB       bool   // is db enable
	DSN        string // dsn connection string
	Pass       string // user identity encryption with cookie
	Pprof      bool   // use profiling in project
}

// Read base config from flags and env.
func InitConfig() *Config {
	config := Config{}
	// Read command line argue.
	startAddress := flag.String("a", "localhost:8080", "start server address and port")
	resultAddress := flag.String("b", "localhost:8080", "answer address and port")
	userAuth := flag.String("s", "mysecret", "User identity encryption with cookie (user_id)")
	tempf := flag.String("f", "", "Location of dump file")
	dsnf := flag.String("d", "", "Data Source Name for DataBase connection")
	pprof := flag.Bool("p", false, "Visualization tool")
	flag.Parse()

	// Check and parse URL.
	startaddr, startport := validators.CheckURL(*startAddress)
	answaddr, answport := validators.CheckURL(*resultAddress)

	config.Address = startaddr + ":" + startport
	config.Response = answaddr + ":" + answport
	zap.S().Infoln("Server address: ", config.Address)

	config.Pass = *userAuth // save cookie pass
	config.Pprof = *pprof   // use pprof in web router

	// Read OS ENV.
	env, exist := os.LookupEnv(("SERVER_ADDRESS"))

	// If env var does not exist  - set def value.
	if exist {
		config.Response = env
		zap.S().Infoln("Set result address from evn SERVER_ADDRESS: ", config.Response)
	} else {
		zap.S().Infoln("Env var SERVER_ADDRESS not found, use default", config.Response)
	}

	// Init Storage.
	dsn, exist := os.LookupEnv(("DATABASE_DSN"))
	// Init storage DB from env variable.
	if exist {
		zap.S().Infoln("Use DataBase Storge. Found location DATABASE_DSN, use: ", dsn)
		config.DSN = dsn
		config.IsDB = true
	} else if *dsnf != "" {
		dsn = *dsnf
		zap.S().Infoln("Use DataBase Storge. Found -d flag, use: ", dsn)
		config.DSN = dsn
		config.IsDB = true
	}

	// Define backup file.
	temp, exist := os.LookupEnv(("FILE_STORAGE_PATH"))
	if exist {
		zap.S().Infoln("Found backup's evn, use file: ", temp)
		config.IsBackup = true
		config.BackupPath = temp
	} else if *tempf != "" {
		zap.S().Infoln("Found backup's flag, use file: ", *tempf)
		config.IsBackup = true
		config.BackupPath = *tempf
	} else {
		config.IsBackup = false
		config.BackupPath = ""
	}

	zap.S().Infoln("Configuration complete")
	return &config
}
