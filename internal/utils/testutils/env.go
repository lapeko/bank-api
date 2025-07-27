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
