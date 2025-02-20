package main

import (
	"net/http"
)

type WebApi struct {
	smtp      SmtpServerConfig
	datastore *DataStore
}

func (web_api *WebApi) Register(mux *http.ServeMux) {

}
