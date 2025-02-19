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
	smtp_server, err := BuildSmtpServerConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("sqlite3", "/data/smtp-sender.v1.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	web_api := WebApi{smtp_server, db}

	mux := http.NewServeMux()
	web_api.Register(mux)

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
		log.Print("using live mode")
		return http.FS(os.DirFS("static"))
	}

	log.Print("using embed mode")
	fsys, err := fs.Sub(embededFiles, "static")
	if err != nil {
		panic(err)
	}

	return http.FS(fsys)
}
