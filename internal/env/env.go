package env

import (
	// "log"
	"os"
	"strconv"
)

func GetString(key string, fallback string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return fallback
}

func GetInt(key string, fallback int) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	number, err := strconv.Atoi(val)

	if err != nil {
		return fallback
	}

	return number
}
