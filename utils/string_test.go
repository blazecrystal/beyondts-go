package utils

import (
	"fmt"
	"testing"
	"time"
)

func TestRandomString(t *testing.T) {
	genRST(make(chan int))
}

func TestConcat(t *testing.T) {
	fmt.Print("**")
	fmt.Print(Concat("aaa", time.Now()))
	fmt.Print("**")
}

func TestX(t *testing.T) {
	ch := make(chan int)
	for i := 0; i < 10; i++ {
		go genRST(ch)
	}
	count := 0
	for {
		count += <-ch
		if count == 10 {
			break
		}
	}
}

func genRST(ch chan int) {
	start := time.Now()
	for i := 0; i < 10; i++ {
		fmt.Println(i, ":", RandomString(10))
	}
	fmt.Println(time.Since(start).Nanoseconds())
	ch <- 1
}
