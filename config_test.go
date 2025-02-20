package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildConfig(t *testing.T) {
	{
		_, err := BuildSmtpServerConfig()
		assert.NotNil(t, err)
	}

	{
		os.Setenv(ConstEnvSmtpServer, "ConstEnvSmtpServer")
		_, err := BuildSmtpServerConfig()
		assert.NotNil(t, err)
	}

	{
		os.Setenv(ConstEnvSmtpPort, "ConstEnvSmtpPort")
		_, err := BuildSmtpServerConfig()
		assert.NotNil(t, err)
	}

	{
		os.Setenv(ConstEnvSmtpAuth, "")
		_, err := BuildSmtpServerConfig()
		assert.NotNil(t, err)
	}

	{
		os.Setenv(ConstEnvSmtpPort, "123")
		_, err := BuildSmtpServerConfig()
		assert.Nil(t, err)
	}

	{
		os.Setenv(ConstEnvSmtpTls, "ConstEnvSmtpTls")
		_, err := BuildSmtpServerConfig()
		assert.NotNil(t, err)
	}

	os.Setenv(ConstEnvSmtpTls, "tls")
	smtpConfig, err := BuildSmtpServerConfig()
	assert.Nil(t, err)

	if assert.NotNil(t, smtpConfig) {
		assert.Equal(t, "ConstEnvSmtpServer", smtpConfig.Server)
		assert.Equal(t, 123, smtpConfig.Port)
		assert.Equal(t, false, smtpConfig.Auth)
		assert.Equal(t, "tls", smtpConfig.Tls)
	}

	{
		os.Setenv(ConstEnvSmtpAuth, "1")
		_, err := BuildSmtpServerConfig()
		assert.NotNil(t, err)
	}

	{
		os.Setenv(ConstEnvSmtpUser, "ConstEnvSmtpUser")
		_, err := BuildSmtpServerConfig()
		assert.NotNil(t, err)
	}

	{
		os.Setenv(ConstEnvSmtpPass, "ConstEnvSmtpPass")
		_, err := BuildSmtpServerConfig()
		assert.Nil(t, err)
	}

	smtpConfig, err = BuildSmtpServerConfig()
	assert.Nil(t, err)
	if assert.NotNil(t, smtpConfig) {
		assert.Equal(t, "ConstEnvSmtpServer", smtpConfig.Server)
		assert.Equal(t, 123, smtpConfig.Port)
		assert.Equal(t, true, smtpConfig.Auth)
		assert.Equal(t, "ConstEnvSmtpUser", smtpConfig.User)
		assert.Equal(t, "ConstEnvSmtpPass", smtpConfig.Pass)
		assert.Equal(t, "tls", smtpConfig.Tls)
	}
}
