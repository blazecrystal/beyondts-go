package utils

import (
	"fmt"
	"testing"
)

func TestGetEnvs(t *testing.T) {
	envs := GetEnvs()
	for k, v := range envs {
		fmt.Println(k, "=", v)
	}
}
