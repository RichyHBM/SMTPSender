package main

import (
	"database/sql"
	"embed"
	"errors"
	"io/fs"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed static
var embededFiles embed.FS

func main() {
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

	useOS := os.Getenv("DEBUG") == "1"
	fs := http.FileServer(getFileSystem(useOS))
	mux.Handle("/", fs)

	log.Println("Server starting up on: ':8080'")

	err = http.ListenAndServe(":8080", mux)
	if errors.Is(err, http.ErrServerClosed) {
		log.Println("Server closed")
	} else if err != nil {
		log.Fatalf("Error starting server: %s\n", err)
	}
}

func getFileSystem(useOS bool) http.FileSystem {
	if useOS {
		log.Print("Using live mode")
		return http.FS(os.DirFS("static"))
	}

	log.Print("Using embed mode")
	fsys, err := fs.Sub(embededFiles, "static")
	if err != nil {
		panic(err)
	}

	return http.FS(fsys)
}
