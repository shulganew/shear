package config

const DefaultHost string = "localhost:8080"

type ConfigShear struct {
	//flag -a
	StartAddress string
	//env var
	ResultAddress string
}

var configapp ConfigShear

func GetConfig() *ConfigShear {

	return &configapp
}
