package crypto

import (
    "crypto/md5"
    "github.com/blazecrystal/beyondts-go/utils"
    "crypto/sha1"
    "crypto/sha256"
    "hash"
    "crypto/sha512"
)

var (
    MD5 = md5.New()
    SHA1 = sha1.New()
    SHA256 = sha256.New()
    SHA512 = sha512.New()
    CHINA_SM3 = NewSM3()
)

// message digest, with given hash algorithm.
// hash should be MD5,SHA1,SHA256,SHA512,CHINA_SM3, which is defined in current file
// src is the bytes to be digested
func Digest(hash hash.Hash, src []byte) []byte {
    hash.Write(src)
    return hash.Sum(nil)
}

// message digest, with given hash algorithm.
// hash should be MD5,SHA1,SHA256,SHA512,CHINA_SM3, which is defined in current file
// src is the string to be digested
// if base64Encoding is true, result string will be encoded with base64, else will be encoded with hex
func DigestString(hash hash.Hash, src string, base64Encoding bool) string {
    tmp := Digest(hash, []byte(src))
    return utils.Bytes2String(tmp, base64Encoding)
}

// message digest, with given hash algorith. after digested, src will be added to result.
// hash should be MD5,SHA1,SHA256,SHA512,CHINA_SM3, which is defined in current file
// src is the bytes to be digested
func DigestX(hash hash.Hash, src []byte) []byte {
    tmp := Digest(hash, src)
    for i := 0; i < len(tmp); i++ {
        tmp[i] = tmp[i] + src[i % len(src)]
    }
    return tmp;
}

// message digest, with given hash algorithm. after digested, src will be added to result.
// hash should be MD5,SHA1,SHA256,SHA512,CHINA_SM3, which is defined in current file
// src is the string to be digested
// if base64Encoding is true, result string will be encoded with base64, else will be encoded with hex
func DigestXString(hash hash.Hash, src string, base64Encoding bool) string {
    tmp := DigestX(hash, []byte(src))
    return utils.Bytes2String(tmp, base64Encoding)
}