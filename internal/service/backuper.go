package service

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/shulganew/shear.git/internal/storage"
	"go.uber.org/zap"
)

type Backup struct {
	File     string
	IsActive bool
}

// {"uuid":"1","short_url":"4rSPg8ap","original_url":"http://yandex.ru"}
func (b Backup) Save(short storage.Short) error {

	data, err := json.Marshal(short)
	zap.S().Infoln("Backup URL: ", string(data))
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

func (b Backup) Load() ([]storage.Short, error) {
	//read
	//var
	file, err := os.OpenFile(b.File, os.O_RDONLY, 0666)
	if err != nil {
		zap.S().Errorln("Error reading backup file")
		return nil, err
	}
	defer file.Close()
	shorts := []storage.Short{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data := scanner.Bytes()
		zap.S().Infoln("Scan bytes", string(data))
		var short storage.Short

		if err := json.Unmarshal(data, &short); err != nil {
			zap.S().Errorln("Error unmarshal data")
			return nil, err
		}
		shorts = append(shorts, short)

	}

	zap.S().Infoln("Load data from file", shorts)
	return shorts, nil
}

func New(file string, isActive bool) *Backup {
	return &Backup{File: file, IsActive: isActive}
}
