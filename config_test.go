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
		os.Setenv(Const_Env_Smtp_Server, "Const_Env_Smtp_Server")
		_, err := BuildSmtpServerConfig()
		assert.NotNil(t, err)
	}

	{
		os.Setenv(Const_Env_Smtp_Port, "Const_Env_Smtp_Port")
		_, err := BuildSmtpServerConfig()
		assert.NotNil(t, err)
	}

	{
		os.Setenv(Const_Env_Smtp_Auth, "")
		_, err := BuildSmtpServerConfig()
		assert.NotNil(t, err)
	}

	{
		os.Setenv(Const_Env_Smtp_Port, "123")
		_, err := BuildSmtpServerConfig()
		assert.Nil(t, err)
	}

	{
		os.Setenv(Const_Env_Smtp_Tls, "Const_Env_Smtp_Tls")
		_, err := BuildSmtpServerConfig()
		assert.NotNil(t, err)
	}

	os.Setenv(Const_Env_Smtp_Tls, "tls")
	smtp_config, err := BuildSmtpServerConfig()
	assert.Nil(t, err)

	if assert.NotNil(t, smtp_config) {
		assert.Equal(t, "Const_Env_Smtp_Server", smtp_config.Server)
		assert.Equal(t, 123, smtp_config.Port)
		assert.Equal(t, false, smtp_config.Auth)
		assert.Equal(t, "tls", smtp_config.Tls)
	}

	{
		os.Setenv(Const_Env_Smtp_Auth, "1")
		_, err := BuildSmtpServerConfig()
		assert.NotNil(t, err)
	}

	{
		os.Setenv(Const_Env_Smtp_User, "Const_Env_Smtp_User")
		_, err := BuildSmtpServerConfig()
		assert.NotNil(t, err)
	}

	{
		os.Setenv(Const_Env_Smtp_Pass, "Const_Env_Smtp_Pass")
		_, err := BuildSmtpServerConfig()
		assert.Nil(t, err)
	}

	smtp_config, err = BuildSmtpServerConfig()
	assert.Nil(t, err)
	if assert.NotNil(t, smtp_config) {
		assert.Equal(t, "Const_Env_Smtp_Server", smtp_config.Server)
		assert.Equal(t, 123, smtp_config.Port)
		assert.Equal(t, true, smtp_config.Auth)
		assert.Equal(t, "Const_Env_Smtp_User", smtp_config.User)
		assert.Equal(t, "Const_Env_Smtp_Pass", smtp_config.Pass)
		assert.Equal(t, "tls", smtp_config.Tls)
	}
}
