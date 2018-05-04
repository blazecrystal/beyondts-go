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

// create aes instance, use given key string.
// in this method, key length is set to 128bit, and iv equals to key
func NewAESInstance(key string) (*AES, error) {
    return NewAESInstance5([]byte(key), KEY_LENGTH_128)
}

// create aes instance with given key string & key length.
// key length should be 128, 256, 512 etc.
// in this method, iv equals to key
func NewAESInstance2(key string, keyLength int) (*AES, error) {
    return NewAESInstance5([]byte(key), keyLength)
}

// create aes instance with given key & iv string.
// in this method, key length is set to 128bit
func NewAESInstance3(key, iv string) (*AES, error) {
    return NewAESInstance6([]byte(key), []byte(iv), KEY_LENGTH_128)
}

// create aes instance with given key & iv string & key length.
func NewAESInstance4(key, iv string, keyLength int) (*AES, error) {
    return NewAESInstance6([]byte(key), []byte(iv), keyLength)
}

// create aes instance with given key bytes & key length.
// in this method, iv equals to key
func NewAESInstance5(key []byte, keyLength int) (*AES, error) {
    key = genKey(key, keyLength)
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    return &AES{block:block, key:key, iv:key, keyLength:keyLength}, err
}

// create aes instance with given key/iv bytes & key length.
func NewAESInstance6(key, iv []byte, keyLength int) (*AES, error) {
    key = genKey(key, keyLength)
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    return &AES{block:block, key:key, iv:iv, keyLength:keyLength}, err
}

// encrypt src bytes.
// just support CBC mode
func (a *AES) Encrypt(src []byte) []byte {
    return encrypt(a.block, src, a.key, a.iv)
}

// encrypt src string.
// just supprt CBC mode
// if base64Encoding is true, result will be encoded with base64, else with hex
func (a *AES) EncryptString(src string, base64Encoding bool) string {
    return utils.Bytes2String(a.Encrypt([]byte(src)), base64Encoding)
}

// decrypt encrypted bytes.
// encrypted should be encrypted in CBC mode
func (a *AES) Decrypt(encrypted []byte) []byte {
    return decrypt(a.block, encrypted, a.key, a.iv)
}

// decrypt encrypted string.
// encrypted should be encrypted in CBC mode
// if base64Encoding is true, means the encrypted string is encoded with base64, else is encoded with hex
func (a *AES) DecryptString(encrypted string, base64Encoding bool) (string, error) {
    tmp, err := utils.String2Bytes(encrypted, base64Encoding)
    if err != nil {
        return "", err
    }
    return string(a.Decrypt(tmp)), err
}