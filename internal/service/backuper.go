// Package service represent a Business layer of app model. Contains:
//
// • shortener - base service for handling of short URLs (brief URLS) and base URLs (origin URLS);
//
// • backuper - service backup config by timer or/and graceful shutdown.
//
// • delete - service for batch aggregation and bulk update for delete user's URLs.
package service

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"time"

	"github.com/shulganew/shear.git/internal/entities"
	"go.uber.org/zap"
)

// Make backup every 10 seconds.
const Timer = 10

// Contain file for backup app data and backups methods.
type Backup struct {
	File string
}

// Backup constructor.
func NewBackup(file string) *Backup {
	return &Backup{File: file}
}

// Write Short entity to file:
//
//	(service.Backup).File
func (b Backup) Save(short entities.Short) error {
	data, err := json.Marshal(short)
	if err != nil {
		zap.S().Error("Error Marshal Backup: ", err)
	}
	file, error := os.OpenFile(b.File, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666) // open file for  data save
	if error != nil {
		return error
	}
	defer file.Close()
	data = append(data, []byte("\n")...) // append next line character
	file.Write(data)
	return nil
}

// Backup all URL data form storage to file:
//
//	(service.Backup).File
func (b Backup) BackupAll(ctx context.Context, storage StorageURL) error {
	// save data fo file
	file, error := os.OpenFile(b.File, os.O_WRONLY|os.O_CREATE, 0666)
	if error != nil {
		return error
	}
	defer file.Close()
	shorts := storage.GetAll(ctx)

	var data []byte
	for _, short := range shorts {

		shortj, err := json.Marshal(short)
		if err != nil {
			zap.S().Error("Error Marshal Backup: ", err)
		}

		// append next line characte
		shortj = append(shortj, []byte("\n")...)
		data = append(data, shortj...)
	}
	zap.S().Infoln("Backup, # of URLs: ", len(shorts))
	file.Write(data)
	return nil
}

// Load data to storage from backup file.
func (b Backup) Load() ([]entities.Short, error) {
	file, err := os.OpenFile(b.File, os.O_RDONLY, 0666)

	if err != nil {
		if os.IsNotExist(err) {
			zap.S().Infoln("Backup file not exist")
			return []entities.Short{}, nil
		}

		zap.S().Errorln("Error reading backup file")
		return nil, err
	}
	defer file.Close()

	shorts := []entities.Short{}
	dec := json.NewDecoder(file)
	for {
		var short entities.Short
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

// Backup all data from storage to file (service.Backup).File during graceful shutdown.
func Shutdown(storage StorageURL, b Backup) {
	// current context doesn't exist, use background context
	b.BackupAll(context.Background(), storage)
}

// Activate backup by timer, backup data every 10 minutes from storage to (service.Backup).File.
func TimeBackup(ctx context.Context, storage StorageURL, b Backup) {
	backup := time.NewTicker(Timer * time.Minute)
	go func() {
		for {
			<-backup.C
			b.BackupAll(ctx, storage)

		}
	}()
}
