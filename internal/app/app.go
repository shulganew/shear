package app

import (
	"log"

	"github.com/shulganew/shear.git/internal/config"
)

// Init app parameters:
// cmd flag -a abd -b
// env
func Init() *config.ConfigShear {

	//init config
	configApp := config.InitConfig()
	log.Println("Config main: ", configApp)

	return configApp
}
