package utils

import (
	"os"
	"strings"
)

func GetEnvs() map[string]string {
	envs := os.Environ()
	envMap := make(map[string]string, len(envs))
	for _, value := range envs {
		equalIndex := strings.Index(value, "=")
		envMap[value[0:equalIndex]] = value[equalIndex+1:]
	}
	return envMap
}
