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
)

func Digest(hash hash.Hash, src []byte) []byte {
    hash.Write(src)
    return hash.Sum(nil)
}

func DigestString(hash hash.Hash, src string, base64Encoding bool) string {
    tmp := Digest(hash, []byte(src))
    return utils.Bytes2String(tmp, base64Encoding)
}

func DigestX(hash hash.Hash, src []byte) []byte {
    tmp := Digest(hash, src)
    for i := 0; i < len(tmp); i++ {
        tmp[i] = tmp[i] + src[i % len(src)]
    }
    return tmp;
}

func DigestXString(hash hash.Hash, src string, base64Encoding bool) string {
    tmp := DigestX(hash, []byte(src))
    return utils.Bytes2String(tmp, base64Encoding)
}