package crypto

import (
    "testing"
    "crypto/elliptic"
    "fmt"
)

func TestEcdsa(t *testing.T) {
    src := "12345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"
    e := NewECDSAInstance()
    err := e.GenKeyPair(elliptic.P256())
    if err != nil {
        fmt.Println(err)
    }
    r, s, err := e.SignString(src, true)
    if err != nil {
        fmt.Println("err")
    }
    fmt.Println("sign : r =", r, "    s =", s)
    v, err := e.VerifyString(src, r, s, true)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println("verified :", v)
}

func TestEcdsa2(t *testing.T) {
    src := "12345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"
    e := NewECDSAInstance()
    err := e.GenKeyPair(elliptic.P256())
    if err != nil {
        fmt.Println(err)
    }
    pubKeyFile := "D:\\workspaces\\workspace-go\\beyondts\\tmp\\keys\\ecdsa\\pubKey.pem";
    priKeyFile := "D:\\workspaces\\workspace-go\\beyondts\\tmp\\keys\\ecdsa\\priKey.pem"
    e.SaveKeyPair(pubKeyFile, priKeyFile)
    e = NewECDSAInstance()
    e.LoadPrivateKey(priKeyFile)
    r, s, err := e.SignString(src, true)
    if err != nil {
        fmt.Println("err")
    }
    fmt.Println("sign : r =", r, "    s =", s)
    e = NewECDSAInstance()
    e.LoadPublicKey(pubKeyFile)
    v, err := e.VerifyString(src, r, s, true)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println("verified :", v)
}