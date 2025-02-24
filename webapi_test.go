package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	smtpmock "github.com/mocktools/go-smtp-mock/v2"
	"github.com/stretchr/testify/assert"
)

func makeSmtpServerConf(port int) *SmtpServerConfig {
	return &SmtpServerConfig{
		"localhost",
		port,
		false,
		"",
		"",
		ConstTlsNone,
		"localhost",
	}
}

func TestSendEmailBadJson(t *testing.T) {
	jsonStr := []byte(`{}`)
	req, err := http.NewRequest("POST", "/test", bytes.NewBuffer(jsonStr))
	assert.NoError(t, err)

	reqRecorder := httptest.NewRecorder()

	api := WebApi{&SmtpServerConfig{}, &DataStore{}}
	api.sendEmail(reqRecorder, req)

	assert.Equal(t, http.StatusBadRequest, reqRecorder.Code)
}

func TestSendEmailBadServer(t *testing.T) {
	jsonStr := []byte(`{"sender": "test@test.test", "recipients": ["test@foo.bar"], "subject": "Test"}`)
	req, err := http.NewRequest("POST", "/test", bytes.NewBuffer(jsonStr))
	assert.NoError(t, err)

	reqRecorder := httptest.NewRecorder()

	smtpServerConf := makeSmtpServerConf(12345)
	api := WebApi{smtpServerConf, &DataStore{}}
	api.sendEmail(reqRecorder, req)

	assert.Equal(t, http.StatusInternalServerError, reqRecorder.Code)
}

func TestSendEmailSuccess(t *testing.T) {
	jsonStr := []byte(`{"sender": "test@test.test", "recipients": ["test@foo.bar"], "subject": "Test"}`)
	req, err := http.NewRequest("POST", "/test", bytes.NewBuffer(jsonStr))
	assert.NoError(t, err)

	reqRecorder := httptest.NewRecorder()

	hostAddress := "localhost"
	server := smtpmock.New(smtpmock.ConfigurationAttr{
		LogToStdout:       false,
		LogServerActivity: false,
		HostAddress:       hostAddress,
	})

	assert.NoError(t, server.Start())

	smtpServerConfig := makeSmtpServerConf(server.PortNumber())

	api := WebApi{smtpServerConfig, &DataStore{}}
	api.sendEmail(reqRecorder, req)

	assert.NoError(t, server.Stop())
	assert.Equal(t, http.StatusOK, reqRecorder.Code)
}
