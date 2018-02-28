package crypto

import (
    "crypto/md5"
    "github.com/blazecrystal/beyondts-go/utils"
    "crypto/sha1"
    "crypto/sha256"
    "hash"
    "crypto/sha512"
)

func MD5(src []byte) []byte {
    return md(md5.New(), src)
}

func MD5String(src string, base64Encoding bool) string {
    tmp := MD5([]byte(src))
    return utils.Bytes2String(tmp, base64Encoding)
}

func SHA1(src []byte) []byte {
    return md(sha1.New(), src)
}

func SHA1String(src string, base64Encoding bool) string {
    tmp := SHA1([]byte(src))
    return utils.Bytes2String(tmp, base64Encoding)
}

func SHA256(src []byte) []byte {
    return md(sha256.New(), src)
}

func SHA256String(src string, base64Encoding bool) string {
    tmp := SHA256([]byte(src))
    return utils.Bytes2String(tmp, base64Encoding)
}

func SHA512(src []byte) []byte {
    return md(sha512.New(), src)
}

func SHA512String(src string, base64Encoding bool) string {
    tmp := SHA512([]byte(src))
    return utils.Bytes2String(tmp, base64Encoding)
}

func md(hash hash.Hash, src []byte) []byte {
    hash.Write(src)
    return hash.Sum(nil)
}