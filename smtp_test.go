package main

import (
	"testing"

	smtpmock "github.com/mocktools/go-smtp-mock/v2"
	"github.com/stretchr/testify/assert"
)

func TestSendEmailBadSender(t *testing.T) {
	assert.Error(t, sendEmail(&SmtpServerConfig{}, "foo", []string{""}, "", ""))
}

func TestSendEmailBadRecipientServer(t *testing.T) {
	assert.Error(t, sendEmail(&SmtpServerConfig{}, "foo@foo.bar", []string{"foo"}, "", ""))
}

func TestSendEmailWrongServer(t *testing.T) {
	hostAddress := "localhost"
	server := smtpmock.New(smtpmock.ConfigurationAttr{
		LogToStdout:       false,
		LogServerActivity: false,
		HostAddress:       hostAddress,
	})
	assert.NoError(t, server.Start())

	smtpConfig := &SmtpServerConfig{
		"localhosta",
		server.PortNumber(),
		false,
		"foo",
		"bar",
		ConstTlsNone,
		"localhost",
	}

	assert.Error(t, sendEmail(smtpConfig, "foo@foo.bar", []string{"foo@foo.bar"}, "", ""))
	assert.NoError(t, server.Stop())
}

func TestSendEmailSuccessfull(t *testing.T) {
	hostAddress := "localhost"
	server := smtpmock.New(smtpmock.ConfigurationAttr{
		LogToStdout:       false,
		LogServerActivity: false,
		HostAddress:       hostAddress,
	})
	assert.NoError(t, server.Start())

	smtpConfig := &SmtpServerConfig{
		hostAddress,
		server.PortNumber(),
		false,
		"foo",
		"bar",
		ConstTlsNone,
		hostAddress,
	}

	assert.NoError(t, sendEmail(smtpConfig, "foo@foo.bar", []string{"foo@foo.bar"}, "", ""))
	assert.NoError(t, server.Stop())
}
