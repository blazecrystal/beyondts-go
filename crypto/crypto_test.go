package crypto

import (
    "testing"
    "fmt"
)

func TestDESCryptString(t *testing.T) {
    src := "123456"
    key := "12345678"
    d, err := NewDESInstance(key)
    if err != nil {
        fmt.Println(err)
    }
    base64 := d.DESEncryptString(src, true)
    hex := d.DESEncryptString(src, false)
    fmt.Println("des enc with base64 :", base64)
    fmt.Println("des enc with hex :", hex)
    src, err = d.DESDecryptString(base64, true)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println("des dec with base64 :", src)
    src, err = d.DESDecryptString(hex, false)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println("des dec with hex :", src)
}

func TestTriDESCryptString(t *testing.T) {
    src := "123456"
    key := "12345678abcdefgh!@#$%^&*"
    d, err := NewTriDESInstance(key)
    base64 := d.TriDESEncryptString(src, true)
    hex := d.TriDESEncryptString(src, false)
    fmt.Println("3des enc with base64 :", base64)
    fmt.Println("3des enc with hex :", hex)
    src, err = d.TriDESDecryptString(base64, true)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println("3des dec with base64 :", src)
    src, err = d.TriDESDecryptString(hex, false)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println("3des dec with hex :", src)
}

func TestAESCryptString(t *testing.T) {
    src := "123456"
    key := "12345678abcdefgh!@#$%^&*"
    //fmt.Println("src :", []byte(src))
    //fmt.Println("key :", []byte(key))
    //tmp, err:=AESEncrypt([]byte(src), []byte(key), KEY_LENGTH_AES128)
    //fmt.Println(tmp)
    a, err := NewAESInstance(key)
    fmt.Println("--------key length : 128-----------")
    base64 := a.EncryptString(src, true)
    hex := a.EncryptString(src, false)
    fmt.Println("aes enc with base64 :", base64)
    fmt.Println("aes enc with hex :", hex)
    src, err = a.DecryptString(base64, true)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println("aes dec with base64 :", src)
    src, err = a.DecryptString(hex, false)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println("aes dec with hex :", src)


    fmt.Println("--------key length : 192-----------")
    a, err = NewAESInstance2(key, KEY_LENGTH_192)
    base64 = a.EncryptString(src, true)
    hex = a.EncryptString(src, false)
    fmt.Println("aes enc with base64 :", base64)
    fmt.Println("aes enc with hex :", hex)
    src, err = a.DecryptString(base64, true)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println("aes dec with base64 :", src)
    src, err = a.DecryptString(hex, false)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println("aes dec with hex :", src)

    fmt.Println("--------key length : 256-----------")
    a, err = NewAESInstance2(key, KEY_LENGTH_256)
    base64 = a.EncryptString(src, true)
    hex = a.EncryptString(src, false)
    fmt.Println("aes enc with base64 :", base64)
    fmt.Println("aes enc with hex :", hex)
    src, err = a.DecryptString(base64, true)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println("aes dec with base64 :", src)
    src, err = a.DecryptString(hex, false)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println("aes dec with hex :", src)

}