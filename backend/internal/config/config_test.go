package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func unsetEnvVars(t *testing.T, keys ...string) {
	t.Helper()
	original := make(map[string]string, len(keys))
	for _, k := range keys {
		original[k] = os.Getenv(k)
		os.Unsetenv(k)
	}
	t.Cleanup(func() {
		for k, v := range original {
			if v == "" {
				os.Unsetenv(k)
			} else {
				os.Setenv(k, v)
			}
		}
	})
}

func TestLoad_DefaultPort(t *testing.T) {
	unsetEnvVars(t, "PORT")

	cfg := Load()

	assert.Equal(t, "8080", cfg.Port)
}

func TestLoad_DefaultDBName(t *testing.T) {
	unsetEnvVars(t, "DB_NAME")

	cfg := Load()

	assert.Equal(t, "open_talk.db", cfg.DBName)
}

func TestLoad_DefaultJWTSecret(t *testing.T) {
	unsetEnvVars(t, "JWT_SECRET")

	cfg := Load()

	assert.Equal(t, "default-secret-change-me", cfg.JWTSecret)
}

func TestLoad_DefaultDBHostIsEmpty(t *testing.T) {
	unsetEnvVars(t, "DB_HOST")

	cfg := Load()

	assert.Equal(t, "", cfg.DBHost)
}

func TestLoad_EnvOverridesPort(t *testing.T) {
	unsetEnvVars(t, "PORT")
	os.Setenv("PORT", "9090")

	cfg := Load()

	assert.Equal(t, "9090", cfg.Port)
}

func TestLoad_EnvOverridesJWTSecret(t *testing.T) {
	unsetEnvVars(t, "JWT_SECRET")
	os.Setenv("JWT_SECRET", "my-super-secret")

	cfg := Load()

	assert.Equal(t, "my-super-secret", cfg.JWTSecret)
}

func TestLoad_EnvOverridesDBName(t *testing.T) {
	unsetEnvVars(t, "DB_NAME")
	os.Setenv("DB_NAME", "custom.db")

	cfg := Load()

	assert.Equal(t, "custom.db", cfg.DBName)
}

func TestLoad_EnvOverridesDBHost(t *testing.T) {
	unsetEnvVars(t, "DB_HOST")
	os.Setenv("DB_HOST", "localhost")

	cfg := Load()

	assert.Equal(t, "localhost", cfg.DBHost)
}

func TestIsSQLite_ReturnsTrueWhenDBHostIsEmpty(t *testing.T) {
	cfg := &Config{DBHost: ""}

	assert.True(t, cfg.IsSQLite())
}

func TestIsSQLite_ReturnsFalseWhenDBHostIsSet(t *testing.T) {
	cfg := &Config{DBHost: "postgres.example.com"}

	assert.False(t, cfg.IsSQLite())
}
