package service

import (
	"bufio"
	"encoding/json"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/shulganew/shear.git/internal/storage"
	"go.uber.org/zap"
)

// make backup every 10 seconds
const Time_Backup = 10

type Backup struct {
	File     string
	IsActive bool
}

func (b Backup) Save(short storage.Short) error {

	data, err := json.Marshal(short)
	//Backup URL:
	if err != nil {
		zap.S().Error("Error Marshal Backup: ", err)
	}
	//save data fo file
	file, error := os.OpenFile(b.File, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if error != nil {
		return error
	}
	defer file.Close()
	//append next line characte
	data = append(data, []byte("\n")...)

	file.Write(data)
	return nil
}

func (b Backup) SaveAll(storage storage.StorageURL) error {

	//save data fo file
	file, error := os.OpenFile(b.File, os.O_WRONLY|os.O_CREATE, 0666)
	if error != nil {
		return error
	}
	defer file.Close()
	shorts := storage.GetAll()

	var data []byte
	for _, short := range shorts {

		shortj, err := json.Marshal(short)
		//Backup URL:
		if err != nil {
			zap.S().Error("Error Marshal Backup: ", err)
		}

		//append next line characte
		shortj = append(shortj, []byte("\n")...)
		data = append(data, shortj...)
	}
	file.Write(data)
	return nil
}

func (b Backup) Load() ([]storage.Short, error) {

	file, err := os.OpenFile(b.File, os.O_RDONLY, 0666)

	if err != nil {
		if os.IsNotExist(err) {
			zap.S().Infoln("Backup file not exist")
			return []storage.Short{}, nil
		}

		zap.S().Errorln("Error reading backup file")
		return nil, err
	}
	defer file.Close()
	shorts := []storage.Short{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data := scanner.Bytes()
		//Scan bytes
		var short storage.Short

		if err := json.Unmarshal(data, &short); err != nil {
			zap.S().Errorln("Error unmarshal data")
			return nil, err
		}
		shorts = append(shorts, short)

	}

	zap.S().Infoln("Load dump from file done. Restore # of elements: ", len(shorts))
	return shorts, nil
}

func Shutdown(storage storage.StorageURL, b Backup) {
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGTERM, syscall.SIGQUIT, os.Interrupt)
	go func() {
		<-exit
		b.SaveAll(storage)
		os.Exit(1)
	}()
}

func TimeBackup(storage storage.StorageURL, b Backup) {

	backup := time.NewTicker(Time_Backup * time.Second)
	go func() {
		for {
			<-backup.C
			b.SaveAll(storage)

		}
	}()

}

func New(file string, isActive bool, storage storage.StorageURL) *Backup {
	return &Backup{File: file, IsActive: isActive}
}
