package config

import (
	"testing"

	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/utils/testutils"
	"github.com/stretchr/testify/require"
)

func clearConfig(envMap map[string]string) {
	testutils.ClearEnv(envMap)
	cfg = nil
}

func TestGet(t *testing.T) {
	envMap := map[string]string{"POSTGRES_URL": testutils.GenRandString(6), "JWT_SECRET_KEY": testutils.GenRandString(6)}

	defer clearConfig(envMap)
	testutils.SetEnv(envMap)

	cfg := Get()

	require.Equal(t, cfg.AppPort, 3000)
	require.Equal(t, cfg.JwtSecretKey, envMap["JWT_SECRET_KEY"])
	require.Equal(t, cfg.PostgresUrl, envMap["POSTGRES_URL"])

	cgf2 := Get()
	require.Equal(t, cgf2, cfg)
}

func TestGet_Failure(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic, but none occurred")
		}
	}()

	Get()
}
