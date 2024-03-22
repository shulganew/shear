// Package responsible for reading config from flag or ENV variables. After reading environment during init application config var locate in Config struct.
package config

import (
	"encoding/json"
	"flag"
	"os"

	"github.com/shulganew/shear.git/internal/entities"
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

	// Set defaults values.
	config.Address = "localhost:8080"
	config.Response = "localhost:8080"
	config.Pass = "mypass"

	// Read command line argue.
	fconf := readFlags()
	// Read ENV.
	econf := readENV()
	// Read JSON config if existed.
	var jconf *entities.ConfJSON
	if econf.JSONPath != nil {
		jconf = readJSONConf(*econf.JSONPath)
		zap.S().Infoln("Use JSON config from (env value): ", *econf.JSONPath)
	} else if fconf.JSONPath != nil {
		jconf = readJSONConf(*fconf.JSONPath)
		zap.S().Infoln("Use JSON config from (flag value): ", *fconf.JSONPath)
	}
	// If JSON config existed, load to main config file.
	if jconf != nil {
		loadJSONConfig(&config, *jconf)
	}

	// Load flag config.
	loadFlagConfig(&config, fconf)

	// Load enviroment config.
	loadENVConfig(&config, econf)

	// Check and parse URL.
	startaddr, startport := validators.CheckURL(config.Address, config.IsSequre)
	answaddr, answport := validators.CheckURL(config.Response, config.IsSequre)
	config.Address = startaddr + ":" + startport
	config.Response = answaddr + ":" + answport
	zap.S().Infoln("Server address: ", config.Address)

	// Init storage DB from env variable.
	if config.DSN != "" {
		zap.S().Infoln("Use Data Base storage: ", config.DSN)
		config.IsDB = true
	} else {
		zap.S().Infoln("Use memory storage.")
	}

	// Define backup file.

	if config.BackupPath != "" {
		zap.S().Infoln("Found backup's path: ", config.BackupPath)
		config.IsBackup = true
	} else {
		zap.S().Infoln("Backup disable")
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

// Load JOSN config data.
func readJSONConf(path string) *entities.ConfJSON {
	f, err := os.Open(path)
	if err != nil {
		zap.S().Infoln("Couldn't open file, use defaults: ", err)
	}
	jsonDecoder := json.NewDecoder(f)
	var jconf entities.ConfJSON
	err = jsonDecoder.Decode(&jconf)
	if err != nil {
		zap.S().Infoln("Couldn't unmarshal file, use defauls: ", err)
	}
	return &jconf
}

// Check if flag passed.
func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

// Read flags to DTO object.
func readFlags() entities.ConfFlag {
	fconf := entities.ConfFlag{}
	startAddress := flag.String("a", "", "start server address and port")
	resultAddress := flag.String("b", "", "answer address and port")
	userAuth := flag.String("x", "mysecret", "User identity encryption with cookie (user_id)")
	tempf := flag.String("f", "", "Location of dump file")
	dsnf := flag.String("d", "", "Data Source Name for DataBase connection")
	pprof := flag.Bool("p", false, "Visualization tool")
	seq := flag.Bool("s", false, "Use sequre connection TLS")
	// Read JSON config
	jsonS := flag.String("c", "", "Path to JSON file with configuration")
	jsonL := flag.String("config", "", "Path to JSON file with configuration")
	flag.Parse()

	if isFlagPassed("a") {
		fconf.Address = startAddress
	}
	if isFlagPassed("b") {
		fconf.Response = resultAddress
	}
	if isFlagPassed("x") {
		fconf.Pass = userAuth
	}
	if isFlagPassed("f") {
		fconf.BackupPath = tempf
	}
	if isFlagPassed("d") {
		fconf.DSN = dsnf
	}
	if isFlagPassed("p") {
		fconf.Pprof = pprof
	}
	if isFlagPassed("s") {
		fconf.IsSequre = seq
	}
	if isFlagPassed("c") {
		fconf.JSONPath = jsonS
	} else if isFlagPassed("config") {
		fconf.JSONPath = jsonL
	}

	return fconf

}

// Read ENV to DTO object.
func readENV() entities.ConfENV {
	econf := entities.ConfENV{}

	sa, exist := os.LookupEnv(("SERVER_ADDRESS"))
	if exist {
		econf.Address = &sa
	}

	bu, exist := os.LookupEnv(("BASE_URL"))
	if exist {
		econf.Response = &bu
	}

	backup, exist := os.LookupEnv(("FILE_STORAGE_PATH"))
	if exist {
		econf.BackupPath = &backup
	}

	dsn, exist := os.LookupEnv(("DATABASE_DSN"))
	if exist {
		econf.DSN = &dsn
	}

	_, exist = os.LookupEnv(("ENABLE_HTTPS"))
	if exist {
		econf.IsSequre = pointBool(true)
	}

	jconf, exist := os.LookupEnv(("CONFIG"))
	if exist {
		econf.JSONPath = &jconf
	}

	return econf
}

// Return pointer to bool value.
func pointBool(b bool) *bool {
	return &b
}

// Load config data from json configuration to the main configuration.
func loadJSONConfig(config *Config, jconf entities.ConfJSON) {
	if jconf.Address != nil {
		config.Address = *jconf.Address
	}
	if jconf.Response != nil {
		config.Response = *jconf.Response
	}
	if jconf.BackupPath != nil {
		config.BackupPath = *jconf.BackupPath
	}
	if jconf.DSN != nil {
		config.DSN = *jconf.DSN
	}
	if jconf.Pass != nil {
		config.Pass = *jconf.Pass
	}
	if jconf.IsBackup != nil {
		config.IsBackup = *jconf.IsBackup
	}
	if jconf.Pprof != nil {
		config.Pprof = *jconf.Pprof
	}
	if jconf.IsSequre != nil {
		config.IsSequre = *jconf.IsSequre
	}
}

// Load config data from flag cmd configuration to the main configuration.
func loadFlagConfig(config *Config, fconf entities.ConfFlag) {
	if fconf.Address != nil {
		config.Address = *fconf.Address
	}
	if fconf.Response != nil {
		config.Response = *fconf.Response
	}
	if fconf.BackupPath != nil {
		config.BackupPath = *fconf.BackupPath
	}
	if fconf.DSN != nil {
		config.DSN = *fconf.DSN
	}
	if fconf.Pass != nil {
		config.Pass = *fconf.Pass
	}
	if fconf.IsBackup != nil {
		config.IsBackup = *fconf.IsBackup
	}
	if fconf.Pprof != nil {
		config.Pprof = *fconf.Pprof
	}
	if fconf.IsSequre != nil {
		config.IsSequre = *fconf.IsSequre
	}
}

// Load config data from ENV configuration to the main configuration.
func loadENVConfig(config *Config, econf entities.ConfENV) {
	if econf.Address != nil {
		config.Address = *econf.Address
	}
	if econf.Response != nil {
		config.Response = *econf.Response
	}
	if econf.BackupPath != nil {
		config.BackupPath = *econf.BackupPath
	}
	if econf.DSN != nil {
		config.DSN = *econf.DSN
	}
	if econf.Pass != nil {
		config.Pass = *econf.Pass
	}
	if econf.IsSequre != nil {
		config.IsSequre = *econf.IsSequre
	}
}
