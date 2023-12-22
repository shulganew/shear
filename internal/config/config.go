package config

import (
	"flag"
	"os"

	"github.com/shulganew/shear.git/internal/web/validators"
	"go.uber.org/zap"
)

const DefaultHost string = "localhost:8080"

type Config struct {
	//flag -a
	Address string
	//env var, or flag -b if env not exist
	Response string

	//Is backup enable
	IsBackup bool

	//File Path for backup
	BackupPath string

	//Is db enable
	IsDB bool

	//dsn connection string
	DSN string

	//User identity encription with cookie
	Pass string
}

func InitConfig() *Config {

	config := Config{}

	//read command line argue
	startAddress := flag.String("a", "localhost:8080", "start server address and port")
	resultAddress := flag.String("b", "localhost:8080", "answer address and port")

	userAuth := flag.String("p", "mysecret", "User identity encription with cookie (user_id)")
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

	//save cookie pass
	config.Pass = *userAuth

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
		config.DSN = dsn
		config.IsDB = true

	} else if *dsnf != "" {
		dsn = *dsnf
		zap.S().Infoln("Use DataBase Storge. Found -d flag, use: ", dsn)
		config.DSN = dsn
		config.IsDB = true

	}

	//define backup file

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

	zap.S().Infoln("Configuration complite")
	return &config
}

func (c *Config) SetConfig(startAddress, resultAddress string) {
	c.Address = startAddress
	c.Response = resultAddress
}

func (c *Config) GetStartAddr() string {
	return c.Address
}

func (c *Config) GetResultAddr() string {
	return c.Response
}
