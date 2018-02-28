package crypto

import (
    "crypto/cipher"
    "crypto/aes"
    "github.com/blazecrystal/beyondts-go/utils"
)

type Aes struct {
    block cipher.Block
    key, iv []byte
    keyLength int
}

func NewAesInstance(key string) (*Aes, error) {
    return NewAesInstance5([]byte(key), KEY_LENGTH_AES128)
}

func NewAesInstance2(key string, keyLength int) (*Aes, error) {
    return NewAesInstance5([]byte(key), keyLength)
}

func NewAesInstance3(key, iv string) (*Aes, error) {
    return NewAesInstance6([]byte(key), []byte(iv), KEY_LENGTH_AES128)
}

func NewAesInstance4(key, iv string, keyLength int) (*Aes, error) {
    return NewAesInstance6([]byte(key), []byte(iv), keyLength)
}

func NewAesInstance5(key []byte, keyLength int) (*Aes, error) {
    key = genAesKey(key, keyLength)
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    return &Aes{block:block, key:key, iv:key, keyLength:keyLength}, err
}

func NewAesInstance6(key, iv []byte, keyLength int) (*Aes, error) {
    key = genAesKey(key, keyLength)
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    return &Aes{block:block, key:key, iv:iv, keyLength:keyLength}, err
}

func (a *Aes) AesEncrypt(src []byte) []byte {
    return encrypt(a.block, src, a.key, a.iv)
}

func (a *Aes) AesEncryptString(src string, base64Encoding bool) string {
    return utils.Bytes2String(a.AesEncrypt([]byte(src)), base64Encoding)
}

func (a *Aes) AesDecrypt(encrypted []byte) []byte {
    return decrypt(a.block, encrypted, a.key, a.iv)
}

func (a *Aes) AesDecryptString(encrypted string, base64Encoding bool) (string, error) {
    tmp, err := utils.String2Bytes(encrypted, base64Encoding)
    if err != nil {
        return "", err
    }
    return string(a.AesDecrypt(tmp)), err
}

func genAesKey(key []byte, keyLength int) []byte {
    if keyLength != KEY_LENGTH_AES128 && keyLength != KEY_LENGTH_AES192 && keyLength != KEY_LENGTH_AES256 {
        keyLength = KEY_LENGTH_AES128
    }
    tmp := keyLength / KEY_LENGTH_DES
    return genBytes(key, tmp)
}