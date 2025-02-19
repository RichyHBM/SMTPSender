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
	smtp_server_host := os.Getenv(Const_Env_Smtp_Server)
	if len(smtp_server_host) == 0 {
		return SmtpServerConfig{}, fmt.Errorf("please supply smtp server via env variable %s", Const_Env_Smtp_Server)
	}

	smtp_port_string := os.Getenv(Const_Env_Smtp_Port)
	if len(smtp_port_string) == 0 {
		return SmtpServerConfig{}, fmt.Errorf("please supply smtp server port via env variable %s", Const_Env_Smtp_Port)
	}
	smtp_port, err := strconv.Atoi(smtp_port_string)
	if err != nil {
		return SmtpServerConfig{}, fmt.Errorf("please supply valid port number for smtp server port via env variable %s", Const_Env_Smtp_Port)
	}

	smtp_auth_string := os.Getenv(Const_Env_Smtp_Auth)

	var smtp_auth bool
	var smtp_user string
	var smtp_pass string

	if smtp_auth_string == "true" || smtp_auth_string == "TRUE" || smtp_auth_string == "1" {
		smtp_auth = true
		smtp_user = os.Getenv(Const_Env_Smtp_User)
		if len(smtp_user) == 0 {
			return SmtpServerConfig{}, fmt.Errorf("please supply smtp server user via env variable %s", Const_Env_Smtp_User)
		}

		smtp_pass = os.Getenv(Const_Env_Smtp_Pass)
		if len(smtp_pass) == 0 {
			return SmtpServerConfig{}, fmt.Errorf("please supply smtp server password via env variable %s", Const_Env_Smtp_Pass)
		}
	} else {
		smtp_auth = false
	}

	tls_mode := os.Getenv(Const_Env_Smtp_Tls)
	if len(tls_mode) != 0 {
		if tls_mode != "none" && tls_mode != "insecure-tls" && tls_mode != "tls" {
			return SmtpServerConfig{}, fmt.Errorf("please supply a tls mode value of: 'none', 'insecure-tls', 'tls', via env variable %s, default if empty is none", Const_Env_Smtp_Tls)
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
