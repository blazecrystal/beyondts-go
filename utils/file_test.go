package utils

import (
	"fmt"
	"testing"
)

func TestListDir(t *testing.T) {
	path := "D:\\workspaces\\workspace-go\\study"
	all, err := ListDir(path, "s", true, false)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, e := range all {
		fmt.Println(e)
	}
}
