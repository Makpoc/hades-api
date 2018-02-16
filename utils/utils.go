package utils

import "os"

// GetEnvPropOrDefault tries to get the property from the environment. If it's not set the function returns the provided default value
func GetEnvPropOrDefault(key, def string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return def
}
