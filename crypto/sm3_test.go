package crypto

import (
    "testing"
    "fmt"
)

func TestSM3Digest(t *testing.T) {
    src := "1234567890中国"
    fmt.Println([]byte(src))
    fmt.Println(string([]byte(src)))
    fmt.Println(SM3DigestString(src, true))
    src = "dddfsl"
    fmt.Println(SM3DigestString(src, true))
    src = "1234567890中国"
    fmt.Println(SM3DigestString(src, true))
}
