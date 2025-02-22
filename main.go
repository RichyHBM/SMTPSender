package main

import (
	"database/sql"
	"embed"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/urfave/negroni"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed static
var embededFiles embed.FS

func main() {
	isDebug := os.Getenv("DEBUG") == "1"

	smtpServer, err := BuildSmtpServerConfig()
	if err != nil {
		log.Fatal(err)
	}

	dbFile := "/data/smtp-sender.v1.db"
	if _, err := os.Stat(dbFile); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if dbFile, err := os.Create(dbFile); err != nil {
				log.Fatal(err)
			} else {
				dbFile.Close()
			}
		} else {
			log.Fatal(err)
		}
	}

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if _, err := db.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		log.Fatal(err)
	}

	datastore, err := MakeDataStore(db)
	if err != nil {
		log.Fatal(err)
	}
	defer datastore.Close()

	mux := http.NewServeMux()
	webApi := WebApi{smtpServer, datastore}
	webApi.Register(mux)

	nux := negroni.New()
	nux.Use(negroni.NewRecovery())
	nux.Use(negroni.NewLogger())
	if domainName := os.Getenv("DOMAIN_NAME"); len(domainName) > 0 {
		nux.Use(MakeEnsureHeaderMiddleware("x-forwarded-host", domainName))
		nux.Use(MakeEnsureHeaderMiddleware("Remote-User", ""))
	}
	nux.Use(negroni.NewStatic(GetHttpFileSystem(isDebug)))
	nux.UseHandler(mux)

	err = http.ListenAndServe(":8080", nux)
	if errors.Is(err, http.ErrServerClosed) {
		log.Println("Server closed")
	} else if err != nil {
		log.Fatalf("Error starting server: %s\n", err)
	}
}
