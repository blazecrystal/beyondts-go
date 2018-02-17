package utils

import (
	"bufio"
	"io"
	"os"
	"testing"
)

func TestPipeTo(t *testing.T) {
	f1, err := os.Open("d:\\1.txt")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	PipeTo(bufio.NewReader(f1), bufio.NewWriter(os.Stdout), io.EOF)
}

func TestPipeTo2(t *testing.T) {
	f1, err := os.Open("d:\\1.txt")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	PipeTo(f1, os.Stdout, 10)
}

func TestPipeTo3(t *testing.T) {
	//endFlag:="123uosud9"
	f1, err := os.Open("d:\\1.txt")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	PipeTo(f1, os.Stdout, nil)
}
