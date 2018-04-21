package crypto

import (
    "testing"
    "fmt"
)

func TestMD(t *testing.T) {
    src := "123456"
    fmt.Println("md5 with base64 :", DigestString(MD5, src, true))
    fmt.Println("md5 with hex :", DigestString(MD5, src, false))
    fmt.Println("sha1 with base64 :", DigestString(SHA1, src, true))
    fmt.Println("sha1 with hex :", DigestString(SHA1, src, false))
    fmt.Println("sha256 with base64 :", DigestString(SHA256, src, true))
    fmt.Println("sha256 with hex :", DigestString(SHA256, src, false))
    fmt.Println("sha512 with base64 :", DigestString(SHA512, src, true))
    fmt.Println("sha512 with hex :", DigestString(SHA512, src, false))
    fmt.Println("--------variant--------")
    fmt.Println("md5 with base64 :", DigestXString(MD5, src, true))
    fmt.Println("md5 with hex :", DigestXString(MD5, src, false))
    fmt.Println("sha1 with base64 :", DigestXString(SHA1, src, true))
    fmt.Println("sha1 with hex :", DigestXString(SHA1, src, false))
    fmt.Println("sha256 with base64 :", DigestXString(SHA256, src, true))
    fmt.Println("sha256 with hex :", DigestXString(SHA256, src, false))
    fmt.Println("sha512 with base64 :", DigestXString(SHA512, src, true))
    fmt.Println("sha512 with hex :", DigestXString(SHA512, src, false))

    fmt.Println("sm3 :", DigestString(CHINA_SM3, src, true))
}