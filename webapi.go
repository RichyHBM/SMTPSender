package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type WebApi struct {
	smtp      *SmtpServerConfig
	datastore *DataStore
}

type SendRequest struct {
	Sender     string   `json:"sender"`
	Recipients []string `json:"recipients"`
	Subject    string   `json:"subject"`
	Body       string   `json:"body"`
}

func (web_api *WebApi) register(mux *http.ServeMux) {
	mux.HandleFunc("/api/send", web_api.sendEmail)
}

func (web_api *WebApi) sendEmail(w http.ResponseWriter, req *http.Request) {
	var request SendRequest

	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Empty body is allowed
	if request.Sender == "" || len(request.Recipients) == 0 || request.Subject == "" {
		http.Error(w, "Missing field in request", http.StatusBadRequest)
		return
	}

	if err := sendEmail(web_api.smtp, request.Sender, request.Recipients, request.Subject, request.Body); err != nil {
		log.Printf("Error sending email: %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
