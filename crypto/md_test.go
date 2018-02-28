package crypto

import (
    "testing"
    "fmt"
)

func TestMD(t *testing.T) {
    src := "123456"
    fmt.Println("md5 with base64 :", MD5String(src, true))
    fmt.Println("md5 with hex :", MD5String(src, false))
    fmt.Println("sha1 with base64 :", SHA1String(src, true))
    fmt.Println("sha1 with hex :", SHA1String(src, false))
    fmt.Println("sha256 with base64 :", SHA256String(src, true))
    fmt.Println("sha256 with hex :", SHA256String(src, false))
    fmt.Println("sha512 with base64 :", SHA512String(src, true))
    fmt.Println("sha512 with hex :", SHA512String(src, false))
}