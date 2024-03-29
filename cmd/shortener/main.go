// Shortener - service for short URL generation.
package main

import (
	"context"
	"database/sql"
	"net/http"

	_ "net/http/pprof"
	"os"
	"sync"

	_ "github.com/shulganew/shear.git/docs" // docs is generated by Swag CLI, you have to import it.
	"github.com/shulganew/shear.git/internal/app"
	"github.com/shulganew/shear.git/internal/config"
	"github.com/shulganew/shear.git/internal/service"
	"github.com/shulganew/shear.git/internal/web/router"
	"go.uber.org/zap"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

// @Title Shortener API
// @Description Shortener service.
// @Version 1.0

// @Contact.email shulganew@gmail.com

// @BasePath /
// @Host localhost:8080
func main() {
	app.Intro(buildVersion, buildDate, buildCommit)
	app.InitLog()
	ctx, cancel := app.InitContext()
	defer cancel()

	conf := config.InitConfig()

	var db *sql.DB
	var err error
	if conf.IsDB {
		db, err = app.InitDB(ctx, conf.DSN)
		if err != nil {
			db = nil
			conf.IsDB = false
			zap.S().Errorln("Can't connect to Database!", err)
		}
		defer db.Close()
	}

	// Use fanIn pattern for storing data from delete requests.
	finalCh := make(chan service.DelBatch, 100)
	defer close(finalCh)

	var waitDel sync.WaitGroup

	// Init application.
	short, backup, del := app.InitApp(ctx, *conf, db, finalCh, &waitDel)

	go func(ctx context.Context, short *service.Shorten, backup *service.Backup, finalCh chan service.DelBatch, wg *sync.WaitGroup) {
		for {
			select {
			case <-ctx.Done():
				zap.S().Infoln("Waiting of update delete...")
				wg.Wait()
				if conf.IsBackup {
					service.Shutdown(short, *backup)
				}
				os.Exit(0)
			case delBatch := <-finalCh:
				short.DeleteBatch(ctx, delBatch)
			}
		}
	}(ctx, short, backup, finalCh, &waitDel)

	// Start web server.
	if err := http.ListenAndServe(conf.Address, router.RouteShear(conf, short, db, del, &waitDel)); err != nil {
		panic(err)
	}
}
