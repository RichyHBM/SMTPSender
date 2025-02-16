package main

import "net/http"

type WebApi struct {
	smtp SmtpServer
}

func (web_api *WebApi) Register(mux *http.ServeMux) {

}
