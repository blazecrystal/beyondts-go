package crypto

import (
    "bytes"
    "crypto/cipher"
)

const (
    KEY_LENGTH = 8
    KEY_LENGTH_24 = KEY_LENGTH * 3
    KEY_LENGTH_128 = KEY_LENGTH * 16
    KEY_LENGTH_192 = KEY_LENGTH * 24
    KEY_LENGTH_256 = KEY_LENGTH * 32
    KEY_LENGTH_1024 = KEY_LENGTH * 128
    KEY_LENGTH_2048 = KEY_LENGTH * 256
)

func encrypt(block cipher.Block, src, key, iv []byte) []byte {
    blockSize := block.BlockSize()
    src = pkcs5Padding(src, blockSize)
    mode := cipher.NewCBCEncrypter(block, genBytes(iv, blockSize))
    crypted := make([]byte, len(src))
    mode.CryptBlocks(crypted, src)
    return crypted
}

func decrypt(block cipher.Block, encrypted, key, iv []byte) []byte {
    mode := cipher.NewCBCDecrypter(block, genBytes(iv, block.BlockSize()))
    src := make([]byte, len(encrypted))
    mode.CryptBlocks(src, encrypted)
    return pkcs5UnPadding(src)
}

func pkcs5Padding(data []byte, blockSize int) []byte {
    padding := blockSize - len(data) % blockSize
    padtext := bytes.Repeat([]byte{byte(padding)}, padding)
    return append(data, padtext...)
}

func pkcs5UnPadding(data []byte) []byte {
    length := len(data)
    // 去掉最后一个字节 unpadding 次
    unpadding := int(data[length - 1])
    return data[:(length - unpadding)]
}

func genBytes(originalBytes []byte, length int) []byte {
    tmp := make([]byte, length)
    if len(originalBytes) < length {
        for i := 0; i < length; i++ {
            tmp[i] = originalBytes[i % len(originalBytes)]
        }
    } else {
        for i := 0; i < length; i++ {
            tmp[i] = originalBytes[i]
        }
    }
    return tmp
}

func genKey(key []byte, keyLength int) []byte {
    if keyLength != KEY_LENGTH_128 && keyLength != KEY_LENGTH_192 && keyLength != KEY_LENGTH_256 {
        keyLength = KEY_LENGTH_128
    }
    tmp := keyLength / KEY_LENGTH
    return genBytes(key, tmp)
}




