package main

import (
	"flag"
	"log"
	"os"
	"syscall"

	"github.com/gophertuts/reminders-cli/server"
	"github.com/gophertuts/reminders-cli/server/repositories"
)

var (
	addrFlag        = flag.String("addr", ":8080", "HTTP server address")
	notifierURIFlag = flag.String("notifier", "http://localhost:9000", "Notifier API URI")
	dbFlag          = flag.String("db", "db.json", "Path to db.json file")
	dbCfgFlag       = flag.String("db-cfg", ".db.config.json", "Path to .db.config.json file")
)

func main() {
	flag.Parse()
	db := repositories.NewDB(*dbFlag, *dbCfgFlag)
	appCfg := server.AppConfig{
		Addr:        *addrFlag,
		NotifierURI: *notifierURIFlag,
		DB:          db,
	}
	application := server.New(appCfg)
	go func() {
		err := application.Start()
		if err != nil {
			log.Fatalf("could not start application: %v", err)
		}
	}()
	server.ListenForSignals([]os.Signal{syscall.SIGINT, syscall.SIGTERM}, application, db)
}
