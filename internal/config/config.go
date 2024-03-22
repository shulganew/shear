// Package responsible for reading config from flag or ENV variables. After reading environment during init application config var locate in Config struct.
package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/shulganew/shear.git/internal/web/validators"
	"go.uber.org/zap"
)

// Struct for store main app config.
type Config struct {
	Address    *string `json:"server_address,omitempty"`    // flag -a
	Response   *string `json:"response_address,omitempty"`  //env var, or flag -b if env not exist
	BackupPath *string `json:"file_storage_path,omitempty"` // file location Path for backup
	DSN        *string `json:"database_dsn,omitempty"`      // dsn connection string
	Pass       *string `json:"pass,omitempty"`              // user identity encryption with cookie
	Backup     *bool   `json:"enable_backup,omitempty"`     // is backup enable
	DB         *bool   // is db enable
	JSONPath   *string // path to JSON config file
	Pprof      *bool   `json:"enable_pprof,omitempty"` // use profiling in project
	IsSeq      *bool   `json:"enable_https,omitempty"` // use https with TLS
}

// Read base config from flags and env.
func NewConfig() *Config {
	config := Config{}

	// Read ENV.
	econf := readENV()

	// Set env values to config.
	loadConfig(&config, econf)

	// Read command line argue.
	fconf := readFlags()

	// Set flag values to config.
	loadConfig(&config, fconf)

	// Read JSON config if existed.
	if config.JSONPath != nil {
		jconf := readJSONConf(*econf.JSONPath)
		loadConfig(&config, *jconf)

	}

	// Set defaults values on empty (nil) config values.
	loadConfig(&config, DefaultConfig())

	// Check and parse URL.
	startaddr, startport := validators.CheckURL(config.GetAddress(), config.IsSequre())
	answaddr, answport := validators.CheckURL(config.GetAddress(), config.IsSequre())
	config.SetAddress(startaddr + ":" + startport)
	config.SetResponse(answaddr + ":" + answport)
	zap.S().Infoln("Server address: ", config.GetAddress())

	// Init storage DB from env variable.
	if config.DSN != nil {
		zap.S().Infoln("Use Data Base storage: ", config.GetDSN())
		config.SetIsDB(true)
	} else {
		zap.S().Infoln("Use memory storage.")
	}

	// Define backup file.
	if config.BackupPath != nil {
		zap.S().Infoln("Found backup's path: ", config.GetBackupPath())
		config.SetIsBackup(true)
	} else {
		zap.S().Infoln("Backup disable")
	}

	zap.S().Infoln("Configuration complete:")
	zap.S().Infoln(config.String())

	return &config
}

// Return suffix http or https depend on type of connection (sequre or not).
func (c Config) GetProtocol() string {
	if *c.IsSeq {
		return "https"
	}
	return "http"
}

// Address from config.
func (c Config) GetAddress() string {
	return *c.Address
}

// Set address to config.
func (c *Config) SetAddress(a string) {
	c.Address = &a
}

// Get response address from config.
func (c Config) GetResponse() string {
	return *c.Response
}

// Set response address to config.
func (c *Config) SetResponse(r string) {
	c.Response = &r
}

// Path to backup file from config.
func (c Config) GetBackupPath() string {
	return *c.BackupPath
}

// Data base DSN from config.
func (c Config) GetDSN() string {
	return *c.DSN
}

// Def pass from config for cookie auth.
func (c Config) GetPass() string {
	return *c.Pass
}

// Is backup enable in config.
func (c Config) IsBackup() bool {
	return *c.Backup
}

// Set backup usage.
func (c *Config) SetIsBackup(b bool) {
	c.Backup = &b
}

// Return true if db use, false - memory use.
func (c Config) IsDB() bool {
	return *c.DB
}

// Set db or memory.
func (c *Config) SetIsDB(b bool) {
	c.DB = &b
}

// Is Ppor is enable.
func (c Config) IsPprof() bool {
	return *c.Pprof
}

// Is https enable.
func (c Config) IsSequre() bool {
	return *c.IsSeq
}

// Stringer interface.
func (c Config) String() string {
	var con strings.Builder
	con.WriteString(fmt.Sprintf("\nAddress: %s \n", c.GetAddress()))
	con.WriteString(fmt.Sprintf("Response: %s \n", c.GetResponse()))
	con.WriteString(fmt.Sprintf("BackupPath: %s \n", c.GetBackupPath()))
	con.WriteString(fmt.Sprintf("DSN: %s \n", c.GetDSN()))
	con.WriteString(fmt.Sprintf("Pass: %s \n", c.GetPass()))
	con.WriteString(fmt.Sprintf("IsBackup: %t \n", c.IsBackup()))
	con.WriteString(fmt.Sprintf("Use DB: %t \n", c.IsDB()))
	con.WriteString(fmt.Sprintf("Pprof: %t \n", c.IsPprof()))
	con.WriteString(fmt.Sprintf("IsSequre: %t \n", c.IsSequre()))
	return con.String()
}

// Load JOSN config data.
func readJSONConf(path string) *Config {
	f, err := os.Open(path)
	if err != nil {
		zap.S().Infoln("Couldn't open file, use defaults: ", err)
	}
	jsonDecoder := json.NewDecoder(f)
	var jconf Config
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

// Read flags to Config object.
func readFlags() Config {
	fconf := Config{}
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
		fconf.IsSeq = seq
	}
	if isFlagPassed("c") {
		fconf.JSONPath = jsonS
	} else if isFlagPassed("config") {
		fconf.JSONPath = jsonL
	}

	return fconf

}

// Read ENV to Config object.
func readENV() Config {
	econf := Config{}

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
		econf.IsSeq = ptBool(true)
	}

	jconf, exist := os.LookupEnv(("CONFIG"))
	if exist {
		econf.JSONPath = &jconf
	}

	return econf
}

// Return config object with preset defaults walues.
func DefaultConfig() Config {
	dconf := Config{}
	// Set defaults values.
	dconf.Address = ptStr("localhost:8080")
	dconf.Response = ptStr("localhost:8080")
	dconf.Pass = ptStr("mypass")
	dconf.DB = ptBool(false)
	dconf.Backup = ptBool(false)
	dconf.DSN = ptStr("postgresql://short:1@localhost/short")
	dconf.Pprof = ptBool(false)
	dconf.IsSeq = ptBool(false)
	return dconf
}

func ptBool(b bool) *bool {
	return &b
}

func ptStr(s string) *string {
	return &s
}

// Assing field from loaded config to main config if values not set in main and existed in loaded.
func loadConfig(main *Config, loaded Config) {
	if main.Address == nil && loaded.Address != nil {
		main.Address = loaded.Address
	}
	if main.Response == nil && loaded.Response != nil {
		main.Response = loaded.Response
	}
	if main.BackupPath == nil && loaded.BackupPath != nil {
		main.BackupPath = loaded.BackupPath
	}
	if main.DSN == nil && loaded.DSN != nil {
		main.DSN = loaded.DSN
	}
	if main.Pass == nil && loaded.Pass != nil {
		main.Pass = loaded.Pass
	}
	if main.Backup == nil && loaded.Backup != nil {
		main.Backup = loaded.Backup
	}
	if main.DB == nil && loaded.DB != nil {
		main.DB = loaded.DB
	}
	if main.JSONPath == nil && loaded.JSONPath != nil {
		main.JSONPath = loaded.JSONPath
	}
	if main.Pprof == nil && loaded.Pprof != nil {
		main.Pprof = loaded.Pprof
	}
	if main.IsSeq == nil && loaded.IsSeq != nil {
		main.IsSeq = loaded.IsSeq
	}
}
