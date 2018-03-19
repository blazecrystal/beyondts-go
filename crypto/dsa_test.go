package crypto

import (
    "testing"
    "crypto/dsa"
    "fmt"
)

func TestDsa(t *testing.T) {
    src := "12345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"
    d := NewDSAInstance()
    d.GenKeyPair(dsa.L3072N256)
    r, s, err := d.SignString(src, true)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println("sign : r =", r, "   s =", s)
    v, err := d.VerifyString(src, r, s, true)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println("verified : ", v)
}

