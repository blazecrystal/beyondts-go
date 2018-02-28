package crypto

import (
    "github.com/blazecrystal/beyondts-go/utils"
    "crypto/des"
    "crypto/cipher"
)

type Des struct {
    block cipher.Block
    key, iv []byte
}

func NewDesInstance(key string) (*Des, error) {
    return NewDesInstance3([]byte(key))
}

func NewDesInstance2(key, iv string) (*Des, error) {
    return NewDesInstance4([]byte(key), []byte(iv))
}

func NewDesInstance3(key []byte) (*Des, error) {
    key = genBytes(key, KEY_LENGTH_DES)
    block, err := des.NewCipher(key)
    if err != nil {
        return nil, err
    }
    return &Des{block:block, key:key, iv:key}, err
}

func NewDesInstance4(key, iv []byte) (*Des, error) {
    key = genBytes(key, KEY_LENGTH_DES)
    block, err := des.NewCipher(key)
    if err != nil {
        return nil, err
    }
    return &Des{block:block, key:key, iv:iv}, err
}

func NewTriDesInstance(key string) (*Des, error) {
    return NewTriDesInstance3([]byte(key))
}

func NewTriDesInstance2(key, iv string) (*Des, error) {
    return NewTriDesInstance4([]byte(key), []byte(iv))
}

func NewTriDesInstance3(key []byte) (*Des, error) {
    key = genBytes(key, KEY_LENGTH_3DES)
    block, err := des.NewTripleDESCipher(key)
    if err != nil {
        return nil, err
    }
    return &Des{block:block, key:key, iv:key}, err
}

func NewTriDesInstance4(key, iv []byte) (*Des, error) {
    key = genBytes(key, KEY_LENGTH_3DES)
    block, err := des.NewTripleDESCipher(key)
    if err != nil {
        return nil, err
    }
    return &Des{block:block, key:key, iv:iv}, err
}

func (d *Des) DesEncrypt(src []byte) []byte {
    return encrypt(d.block, src, d.key, d.iv)
}

func (d *Des) DesEncryptString(src string, base64Encoding bool) string {
    tmp := d.DesEncrypt([]byte(src))
    return utils.Bytes2String(tmp, base64Encoding)
}

func (d *Des) DesDecrypt(encrypted []byte) []byte {
    return decrypt(d.block, encrypted, d.key, d.iv)
}

func (d *Des) DesDecryptString(encrypted string, base64Encoding bool) (string, error) {
    tmp, err := utils.String2Bytes(encrypted, base64Encoding)
    if err != nil {
        return "", err
    }
    return string(d.DesDecrypt(tmp)), err
}

func (d *Des) TriDesEncrypt(src []byte) []byte {
    return encrypt(d.block, src, d.key, d.iv)
}

func (d *Des) TriDesEncryptString(src string, base64Encoding bool) string {
    tmp := d.TriDesEncrypt([]byte(src))
    return utils.Bytes2String(tmp, base64Encoding)
}

func (d *Des) TriDesDecrypt(encrypted []byte) []byte {
    return decrypt(d.block, encrypted, d.key, d.iv)
}

func (d *Des) TriDesDecryptString(encrypted string, base64Encoding bool) (string, error) {
    tmp, err := utils.String2Bytes(encrypted, base64Encoding)
    if err != nil {
        return "", err
    }
    return string(d.TriDesDecrypt(tmp)), err
}