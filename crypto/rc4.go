package crypto

import (
    "crypto/cipher"
    "crypto/aes"
    "github.com/blazecrystal/beyondts-go/utils"
)

type RC4 struct {
    block cipher.Block
    key, iv []byte
    keyLength int
}

func NewRC4Instance(key string) (*RC4, error) {
    return NewRC4Instance5([]byte(key), KEY_LENGTH_128)
}

func NewRC4Instance2(key string, keyLength int) (*RC4, error) {
    return NewRC4Instance5([]byte(key), keyLength)
}

func NewRC4Instance3(key, iv string) (*RC4, error) {
    return NewRC4Instance6([]byte(key), []byte(iv), KEY_LENGTH_128)
}

func NewRC4Instance4(key, iv string, keyLength int) (*RC4, error) {
    return NewRC4Instance6([]byte(key), []byte(iv), keyLength)
}

func NewRC4Instance5(key []byte, keyLength int) (*RC4, error) {
    key = genKey(key, keyLength)
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    return &RC4{block:block, key:key, iv:key, keyLength:keyLength}, err
}

func NewRC4Instance6(key, iv []byte, keyLength int) (*RC4, error) {
    key = genKey(key, keyLength)
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    return &RC4{block:block, key:key, iv:iv, keyLength:keyLength}, err
}

func (a *RC4) Encrypt(src []byte) []byte {
    return encrypt(a.block, src, a.key, a.iv)
}

func (a *RC4) EncryptString(src string, base64Encoding bool) string {
    return utils.Bytes2String(a.Encrypt([]byte(src)), base64Encoding)
}

func (a *RC4) Decrypt(encrypted []byte) []byte {
    return decrypt(a.block, encrypted, a.key, a.iv)
}

func (a *RC4) DecryptString(encrypted string, base64Encoding bool) (string, error) {
    tmp, err := utils.String2Bytes(encrypted, base64Encoding)
    if err != nil {
        return "", err
    }
    return string(a.Decrypt(tmp)), err
}
