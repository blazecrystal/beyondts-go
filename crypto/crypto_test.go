package crypto

import (
    "testing"
    "fmt"
)

func TestDesCryptString(t *testing.T) {
    src := "123456"
    key := "12345678"
    d, err := NewDesInstance(key)
    if err != nil {
        fmt.Println(err)
    }
    base64 := d.DesEncryptString(src, true)
    hex := d.DesEncryptString(src, false)
    fmt.Println("des enc with base64 :", base64)
    fmt.Println("des enc with hex :", hex)
    src, err = d.DesDecryptString(base64, true)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println("des dec with base64 :", src)
    src, err = d.DesDecryptString(hex, false)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println("des dec with hex :", src)
}

func TestTriDesCryptString(t *testing.T) {
    src := "123456"
    key := "12345678abcdefgh!@#$%^&*"
    d, err := NewTriDesInstance(key)
    base64 := d.TriDesEncryptString(src, true)
    hex := d.TriDesEncryptString(src, false)
    fmt.Println("3des enc with base64 :", base64)
    fmt.Println("3des enc with hex :", hex)
    src, err = d.TriDesDecryptString(base64, true)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println("3des dec with base64 :", src)
    src, err = d.TriDesDecryptString(hex, false)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println("3des dec with hex :", src)
}

func TestAesCryptString(t *testing.T) {
    src := "123456"
    key := "12345678abcdefgh!@#$%^&*"
    //fmt.Println("src :", []byte(src))
    //fmt.Println("key :", []byte(key))
    //tmp, err:=AesEncrypt([]byte(src), []byte(key), KEY_LENGTH_AES128)
    //fmt.Println(tmp)
    a, err := NewAesInstance(key)
    fmt.Println("--------key length : 128-----------")
    base64 := a.AesEncryptString(src, true)
    hex := a.AesEncryptString(src, false)
    fmt.Println("aes enc with base64 :", base64)
    fmt.Println("aes enc with hex :", hex)
    src, err = a.AesDecryptString(base64, true)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println("aes dec with base64 :", src)
    src, err = a.AesDecryptString(hex, false)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println("aes dec with hex :", src)


    fmt.Println("--------key length : 192-----------")
    a, err = NewAesInstance2(key, KEY_LENGTH_AES192)
    base64 = a.AesEncryptString(src, true)
    hex = a.AesEncryptString(src, false)
    fmt.Println("aes enc with base64 :", base64)
    fmt.Println("aes enc with hex :", hex)
    src, err = a.AesDecryptString(base64, true)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println("aes dec with base64 :", src)
    src, err = a.AesDecryptString(hex, false)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println("aes dec with hex :", src)

    fmt.Println("--------key length : 256-----------")
    a, err = NewAesInstance2(key, KEY_LENGTH_AES256)
    base64 = a.AesEncryptString(src, true)
    hex = a.AesEncryptString(src, false)
    fmt.Println("aes enc with base64 :", base64)
    fmt.Println("aes enc with hex :", hex)
    src, err = a.AesDecryptString(base64, true)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println("aes dec with base64 :", src)
    src, err = a.AesDecryptString(hex, false)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println("aes dec with hex :", src)

}