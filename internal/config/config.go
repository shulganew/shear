// Package responsible for reading config from flag or ENV variables. After reading environment during init application config var locate in Config struct.
package config

import (
	"flag"
	"os"

	"github.com/shulganew/shear.git/internal/web/validators"
	"go.uber.org/zap"
)

// Struct for store main app config.
type Config struct {
	Address    string // flag -a
	Response   string //env var, or flag -b if env not exist
	BackupPath string // file location Path for backup
	DSN        string // dsn connection string
	Pass       string // user identity encryption with cookie
	IsBackup   bool   // is backup enable
	IsDB       bool   // is db enable
	Pprof      bool   // use profiling in project
	IsSequre   bool   // use https with TLS
}

// Read base config from flags and env.
func NewConfig() *Config {
	config := Config{}
	// Read command line argue.
	startAddress := flag.String("a", "", "start server address and port")
	resultAddress := flag.String("b", "", "answer address and port")
	userAuth := flag.String("x", "mysecret", "User identity encryption with cookie (user_id)")
	tempf := flag.String("f", "", "Location of dump file")
	dsnf := flag.String("d", "", "Data Source Name for DataBase connection")
	pprof := flag.Bool("p", false, "Visualization tool")
	seq := flag.Bool("s", false, "Use sequre connection TLS")
	flag.Parse()

	// SSL enable check.
	config.IsSequre = *seq

	// Read  ENV.
	_, exist := os.LookupEnv(("ENABLE_HTTPS"))
	if exist {
		config.IsSequre = true
	}

	// Check and parse URL.
	startaddr, startport, isDefS := validators.CheckURL(*startAddress, config.IsSequre)
	answaddr, answport, isDefR := validators.CheckURL(*resultAddress, config.IsSequre)

	if isDefS {
		zap.S().Infoln("Use default start address: ", startAddress, startport)
	}
	if isDefR {
		zap.S().Infoln("Use default result address: ", answaddr, answport)
	}

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

// Return suffix http or https depend on type of connection (sequre or not).
func (c Config) GetProtocol() string {
	if c.IsSequre {
		return "https"
	}
	return "http"
}
