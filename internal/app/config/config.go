package config

type Config struct {
	//flag -a
	startAddress string
	//flag -b
	resultAddress string
}

func NewConfig(startAddress string, resultAddress string) *Config {

	return &Config{startAddress: startAddress, resultAddress: resultAddress}
}
