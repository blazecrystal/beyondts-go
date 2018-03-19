package crypto

import (
    "testing"
    "fmt"
    "strings"
)

func TestSaveKeyPair(t *testing.T) {
    r := NewRSAInstance()
    err := r.GenKeyPair(KEY_LENGTH_1024)
    if err != nil {
        fmt.Println(err)
        panic(err)
    }
    err = r.SaveKeyPair("d:\\pub.key", "d:\\pri.key")
    if err != nil {
        fmt.Println(err)
        panic(err)
    }
}

func TestCryption(t *testing.T) {
    src := "123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"
    r := NewRSAInstance()
    err := r.GenKeyPair(KEY_LENGTH_1024)
    if err != nil {
        fmt.Println(err)
        panic(err)
    }
    r.SaveKeyPair("d:\\pubKey.pem", "d:\\priKey.pem")
    base64, err := r.EncryptString(src, true)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println("encrypted with base64 :", base64)
    hex, err := r.EncryptString(src, false)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println("encrypted with hex :", hex)

    src, err = r.DecryptString(base64, true)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println("decrypted with base64 :", src)
    hex, err = r.DecryptString(hex, false)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println("decrypted with hex :", src)
    sign, err := r.SignString(src, true)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println("sign with base64 :", sign)
    err = r.VerifyString(src, sign, true)
    fmt.Println("verified with base64 :", err == nil)
    fmt.Println("==========================")
    r = NewRSAInstance()
    err = r.LoadPublicKey("d:\\pubKey.pem")
    if err != nil {
        fmt.Println(err)
    }
    err = r.LoadPrivateKey("d:\\priKey.pem")
    if err != nil {
        fmt.Println(err)
    }
    base642, err := r.EncryptString(src, true)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println("encrypted from saved pubkey is correct :", err == nil, "  ", strings.EqualFold(base642, base64))
    src2, err := r.DecryptString(base642, true)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println("decrypted from saved prikey is correct :", strings.EqualFold(src, src2))
    fmt.Println("=============key length : 2048==============")
    r = NewRSAInstance()
    r.GenKeyPair(KEY_LENGTH_2048)
    base64, err = r.EncryptString(src, true)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println("encrypted with base64 :", base64)
    src, err = r.DecryptString(base64, true)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println("decrypted with base64 :", src)
}

func TestSignJAVAGO(t *testing.T) {
    //r := NewRsaInstance()
    //r.L
}