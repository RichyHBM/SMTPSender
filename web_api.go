package main

import (
	"database/sql"
	"net/http"
)

type WebApi struct {
	smtp SmtpServerConfig
	db   *sql.DB
}

func (web_api *WebApi) Register(mux *http.ServeMux) {

}
