package crypto

import (
    "github.com/blazecrystal/beyondts-go/utils"
    "crypto/des"
    "crypto/cipher"
)

type DES struct {
    block cipher.Block
    key, iv []byte
}

func NewDESInstance(key string) (*DES, error) {
    return NewDESInstance3([]byte(key))
}

func NewDESInstance2(key, iv string) (*DES, error) {
    return NewDESInstance4([]byte(key), []byte(iv))
}

func NewDESInstance3(key []byte) (*DES, error) {
    key = genBytes(key, KEY_LENGTH)
    block, err := des.NewCipher(key)
    if err != nil {
        return nil, err
    }
    return &DES{block:block, key:key, iv:key}, err
}

func NewDESInstance4(key, iv []byte) (*DES, error) {
    key = genBytes(key, KEY_LENGTH)
    block, err := des.NewCipher(key)
    if err != nil {
        return nil, err
    }
    return &DES{block:block, key:key, iv:iv}, err
}

func NewTriDESInstance(key string) (*DES, error) {
    return NewTriDESInstance3([]byte(key))
}

func NewTriDESInstance2(key, iv string) (*DES, error) {
    return NewTriDESInstance4([]byte(key), []byte(iv))
}

func NewTriDESInstance3(key []byte) (*DES, error) {
    key = genBytes(key, KEY_LENGTH_24)
    block, err := des.NewTripleDESCipher(key)
    if err != nil {
        return nil, err
    }
    return &DES{block:block, key:key, iv:key}, err
}

func NewTriDESInstance4(key, iv []byte) (*DES, error) {
    key = genBytes(key, KEY_LENGTH_24)
    block, err := des.NewTripleDESCipher(key)
    if err != nil {
        return nil, err
    }
    return &DES{block:block, key:key, iv:iv}, err
}

func (d *DES) DESEncrypt(src []byte) []byte {
    return encrypt(d.block, src, d.key, d.iv)
}

func (d *DES) DESEncryptString(src string, base64Encoding bool) string {
    tmp := d.DESEncrypt([]byte(src))
    return utils.Bytes2String(tmp, base64Encoding)
}

func (d *DES) DESDecrypt(encrypted []byte) []byte {
    return decrypt(d.block, encrypted, d.key, d.iv)
}

func (d *DES) DESDecryptString(encrypted string, base64Encoding bool) (string, error) {
    tmp, err := utils.String2Bytes(encrypted, base64Encoding)
    if err != nil {
        return "", err
    }
    return string(d.DESDecrypt(tmp)), err
}

func (d *DES) TriDESEncrypt(src []byte) []byte {
    return encrypt(d.block, src, d.key, d.iv)
}

func (d *DES) TriDESEncryptString(src string, base64Encoding bool) string {
    tmp := d.TriDESEncrypt([]byte(src))
    return utils.Bytes2String(tmp, base64Encoding)
}

func (d *DES) TriDESDecrypt(encrypted []byte) []byte {
    return decrypt(d.block, encrypted, d.key, d.iv)
}

func (d *DES) TriDESDecryptString(encrypted string, base64Encoding bool) (string, error) {
    tmp, err := utils.String2Bytes(encrypted, base64Encoding)
    if err != nil {
        return "", err
    }
    return string(d.TriDESDecrypt(tmp)), err
}