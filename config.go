package main

import (
	"fmt"
	"os"
	"strconv"
)

const (
	ConstEnvSmtpServer = "SMTP_SERVER"
	ConstEnvSmtpPort = "SMTP_PORT"
	ConstEnvSmtpAuth = "SKIP_AUTH"
	ConstEnvSmtpUser = "SMTP_USER"
	ConstEnvSmtpPass = "SMTP_PASS"
	ConstEnvSmtpTls = "TLS_MODE"
)

type SmtpServerConfig struct {
	Server string
	Port   int
	Auth   bool
	User   string
	Pass   string
	Tls    string
}

func BuildSmtpServerConfig() (SmtpServerConfig, error) {
	smtpServerHost := os.Getenv(ConstEnvSmtpServer)
	if len(smtpServerHost) == 0 {
		return SmtpServerConfig{}, fmt.Errorf("please supply smtp server via env variable %s", ConstEnvSmtpServer)
	}

	smtpPortString := os.Getenv(ConstEnvSmtpPort)
	if len(smtpPortString) == 0 {
		return SmtpServerConfig{}, fmt.Errorf("please supply smtp server port via env variable %s", ConstEnvSmtpPort)
	}
	smtpPort, err := strconv.Atoi(smtpPortString)
	if err != nil {
		return SmtpServerConfig{}, fmt.Errorf("please supply valid port number for smtp server port via env variable %s", ConstEnvSmtpPort)
	}

	smtpAuthString := os.Getenv(ConstEnvSmtpAuth)

	var smtpAuth bool
	var smtpUser string
	var smtpPass string

	if smtpAuthString == "true" || smtpAuthString == "TRUE" || smtpAuthString == "1" {
		smtpAuth = true
		smtpUser = os.Getenv(ConstEnvSmtpUser)
		if len(smtpUser) == 0 {
			return SmtpServerConfig{}, fmt.Errorf("please supply smtp server user via env variable %s", ConstEnvSmtpUser)
		}

		smtpPass = os.Getenv(ConstEnvSmtpPass)
		if len(smtpPass) == 0 {
			return SmtpServerConfig{}, fmt.Errorf("please supply smtp server password via env variable %s", ConstEnvSmtpPass)
		}
	} else {
		smtpAuth = false
	}

	tlsMode := os.Getenv(ConstEnvSmtpTls)
	if len(tlsMode) != 0 {
		if tlsMode != "none" && tlsMode != "insecure-tls" && tlsMode != "tls" {
			return SmtpServerConfig{}, fmt.Errorf("please supply a tls mode value of: 'none', 'insecure-tls', 'tls', via env variable %s, default if empty is none", ConstEnvSmtpTls)
		}
	}

	return SmtpServerConfig{
		smtpServerHost,
		smtpPort,
		smtpAuth,
		smtpUser,
		smtpPass,
		tlsMode,
	}, nil
}
