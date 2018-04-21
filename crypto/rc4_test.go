package crypto

import (
    "testing"
    "fmt"
    "strings"
)

func TestRC4(t *testing.T) {
    src := "123456"
    key := "1234567890!@#$%^&*"
    r, err := NewRC4Instance2(key, KEY_LENGTH_256)
    if err != nil {
        fmt.Println(err)
    }
    enc := r.EncryptString(src, true)
    fmt.Println("encrypted :", enc)
    dec, err := r.DecryptString(enc, true)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println("decrypted :", strings.EqualFold(dec, src))

    r, err = NewRC4Instance3(key, key)
    if err != nil {
        fmt.Println(err)
    }
    enc = r.EncryptString(src, true)
    fmt.Println("encrypted :", enc)
    dec, err = r.DecryptString(enc, true)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println("decrypted :", strings.EqualFold(dec, src))
}
