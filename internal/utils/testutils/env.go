package testutils

import "os"

func SetEnv(keyValueMap map[string]string) {
	for key, value := range keyValueMap {
		os.Setenv(key, value)
	}
}

func ClearEnv(keyValueMap map[string]string) {
	for key := range keyValueMap {
		os.Unsetenv(key)
	}
}

func MockEnv() (clearFn func()) {
	os.Setenv("JWT_SECRET_KEY", GenRandString(GenRandIntInRange(4, 10)))
	os.Setenv("POSTGRES_URL", GenRandString(GenRandIntInRange(4, 10)))

	return func() {
		os.Unsetenv("JWT_SECRET_KEY")
		os.Unsetenv("POSTGRES_URL")
	}
}
