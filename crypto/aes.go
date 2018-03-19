package crypto

import (
    "crypto/cipher"
    "crypto/aes"
    "github.com/blazecrystal/beyondts-go/utils"
)

type AES struct {
    block cipher.Block
    key, iv []byte
    keyLength int
}

func NewAESInstance(key string) (*AES, error) {
    return NewAESInstance5([]byte(key), KEY_LENGTH_128)
}

func NewAESInstance2(key string, keyLength int) (*AES, error) {
    return NewAESInstance5([]byte(key), keyLength)
}

func NewAESInstance3(key, iv string) (*AES, error) {
    return NewAESInstance6([]byte(key), []byte(iv), KEY_LENGTH_128)
}

func NewAESInstance4(key, iv string, keyLength int) (*AES, error) {
    return NewAESInstance6([]byte(key), []byte(iv), keyLength)
}

func NewAESInstance5(key []byte, keyLength int) (*AES, error) {
    key = genKey(key, keyLength)
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    return &AES{block:block, key:key, iv:key, keyLength:keyLength}, err
}

func NewAESInstance6(key, iv []byte, keyLength int) (*AES, error) {
    key = genKey(key, keyLength)
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    return &AES{block:block, key:key, iv:iv, keyLength:keyLength}, err
}

func (a *AES) Encrypt(src []byte) []byte {
    return encrypt(a.block, src, a.key, a.iv)
}

func (a *AES) EncryptString(src string, base64Encoding bool) string {
    return utils.Bytes2String(a.Encrypt([]byte(src)), base64Encoding)
}

func (a *AES) Decrypt(encrypted []byte) []byte {
    return decrypt(a.block, encrypted, a.key, a.iv)
}

func (a *AES) DecryptString(encrypted string, base64Encoding bool) (string, error) {
    tmp, err := utils.String2Bytes(encrypted, base64Encoding)
    if err != nil {
        return "", err
    }
    return string(a.Decrypt(tmp)), err
}