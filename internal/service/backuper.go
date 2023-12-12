package service

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"time"

	"github.com/shulganew/shear.git/internal/storage"
	"go.uber.org/zap"
)

// make backup every 10 seconds
const Timer = 10

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

func (b Backup) BackupAll(ctx context.Context, storage storage.StorageURL) error {

	//save data fo file
	file, error := os.OpenFile(b.File, os.O_WRONLY|os.O_CREATE, 0666)
	if error != nil {
		return error
	}
	defer file.Close()
	shorts := storage.GetAll(ctx)

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
	zap.S().Infoln("Backup, # of URLs: ", len(shorts))
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
	dec := json.NewDecoder(file)
	for {
		var short storage.Short
		if err := dec.Decode(&short); err == io.EOF {
			break
		} else if err != nil {
			zap.S().Errorln("Error unmarshal data, check validation of backup file", err)
			panic(err)
		}
		shorts = append(shorts, short)
	}
	zap.S().Infoln("Load dump from file done. Restore # of elements: ", len(shorts))

	return shorts, nil
}

func Shutdown(ctx context.Context, storage storage.StorageURL, b Backup) {
	go func() {
		<-ctx.Done()
		//current context doesn't exist, use background context
		b.BackupAll(context.Background(), storage)
		os.Exit(0)
	}()
}

func TimeBackup(ctx context.Context, storage storage.StorageURL, b Backup) {

	backup := time.NewTicker(Timer * time.Minute)
	go func() {
		for {
			<-backup.C
			b.BackupAll(ctx, storage)

		}
	}()

}

func NewBackup(file string, isActive bool) *Backup {
	return &Backup{File: file, IsActive: isActive}
}
