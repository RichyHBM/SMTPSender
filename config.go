package main

import (
	"fmt"
	"os"
	"strconv"
)

type SmtpServerConfig struct {
	Server string
	Port   int
	Auth   bool
	User   string
	Pass   string
	Tls    string
}

type Mail struct {
	From    string
	To      []string
	Subject string
	Body    string
}

func BuildSmtpServerConfig() (SmtpServerConfig, error) {
	smtp_server_host := os.Getenv("SMTP_SERVER")
	if len(smtp_server_host) == 0 {
		return SmtpServerConfig{}, fmt.Errorf("please supply smtp server via env variable SMTP_SERVER")
	}

	smtp_port_string := os.Getenv("SMTP_PORT")
	if len(smtp_port_string) == 0 {
		return SmtpServerConfig{}, fmt.Errorf("please supply smtp server port via env variable SMTP_PORT")
	}
	smtp_port, err := strconv.Atoi(smtp_port_string)
	if err != nil {
		return SmtpServerConfig{}, fmt.Errorf("please supply valid port number for smtp server port via env variable SMTP_PORT")
	}

	smtp_auth_string := os.Getenv("SKIP_AUTH")
	smtp_auth := true
	if len(smtp_auth_string) != 0 {
		smtp_auth = false
	}

	var smtp_user string
	var smtp_pass string

	if smtp_auth {
		smtp_user := os.Getenv("SMTP_USER")
		if len(smtp_user) == 0 {
			return SmtpServerConfig{}, fmt.Errorf("please supply smtp server user via env variable SMTP_USER")
		}

		smtp_pass := os.Getenv("SMTP_PASS")
		if len(smtp_pass) == 0 {
			return SmtpServerConfig{}, fmt.Errorf("please supply smtp server password via env variable SMTP_PASS")
		}
	}

	tls_mode := os.Getenv("TLS_MODE")
	if len(tls_mode) != 0 {
		if tls_mode != "none" && tls_mode != "insecure-tls" && tls_mode != "tls" {
			return SmtpServerConfig{}, fmt.Errorf("please supply a tls mode value of: 'none', 'insecure-tls', 'tls', via env variable TLS_MODE, default if empty is none")
		}
	}

	return SmtpServerConfig{
		smtp_server_host,
		smtp_port,
		smtp_auth,
		smtp_user,
		smtp_pass,
		tls_mode,
	}, nil
}
