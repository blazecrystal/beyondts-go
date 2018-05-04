package utils

import (
    "testing"
    "fmt"
)

func TestRemoveFromSlice(t *testing.T) {
    s:=make([]int, 5)
    for i:=0;i<len(s);i++ {
        s[i] = i
    }
    fmt.Println(s)
    fmt.Printf("&s : %p\n", &s)
    RemoveFromSlice(s, 3)
    fmt.Printf("&s : %p\n", &s)
    fmt.Println(s)
}