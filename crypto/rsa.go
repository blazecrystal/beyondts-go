package crypto

import (
    "crypto/rsa"
    "encoding/pem"
    "crypto/rand"
    "crypto/x509"
    "os"
    "io/ioutil"
    "errors"
    "github.com/blazecrystal/beyondts-go/utils"
    "crypto"
    "bytes"
)

const (
    KEY_LENGTH_RSA1024 = 1024
    KEY_LENGTH_RSA2048 = 2048

    MAX_ENCRYPT_BLOCK = 117
    MAX_DECRYPT_BLOCK = 128
)

type Rsa struct {
    priKey *rsa.PrivateKey
    pubKey *rsa.PublicKey
}

func NewRsaInstance() *Rsa {
    return &Rsa{}
}

func (r *Rsa) GenKeyPair(keyLength int) error {
    priKey, err := rsa.GenerateKey(rand.Reader, keyLength)
    if err != nil {
        return err
    }
    r.priKey = priKey
    r.pubKey = &priKey.PublicKey
    return err
}

func (r *Rsa) EncryptString(src string, base64Encoding bool) (string, error) {
    tmp, err := r.Encrypt([]byte(src))
    if err != nil {
        return "", err
    }
    return utils.Bytes2String(tmp, base64Encoding), err
}

func (r *Rsa) DecryptString(encrypted string, base64Encoding bool) (string, error) {
    tmp, err := utils.String2Bytes(encrypted, base64Encoding)
    if err != nil {
        return "", err
    }
    tmp, err = r.Decrypt(tmp)
    return string(tmp), err
}

func (r *Rsa) Encrypt(src []byte) ([]byte, error) {
    buf := bytes.Buffer{}
    maxBlockSize := len(r.pubKey.N.Bytes()) - 11
    times := len(src) / maxBlockSize
    var err error
    var tmp []byte
    for i := 0; i < times; i++ {
        tmp, err = rsa.EncryptPKCS1v15(rand.Reader, r.pubKey, src[i * maxBlockSize : (i + 1) * maxBlockSize])
        if err != nil {
            return nil, err
        }
        buf.Write(tmp)
    }
    last := len(src) / maxBlockSize
    if last > 0 {
        tmp, err = rsa.EncryptPKCS1v15(rand.Reader, r.pubKey, src[times * maxBlockSize : ])
        if err != nil {
            return nil ,err
        }
        buf.Write(tmp)
    }
    return buf.Bytes(), err
}

func (r *Rsa) Decrypt(encrypted []byte) ([]byte, error) {
    buf := bytes.Buffer{}
    maxBlockSize := len(r.pubKey.N.Bytes())
    times := len(encrypted) / maxBlockSize
    var err error
    var tmp []byte
    for i := 0; i < times; i++ {
        tmp, err = rsa.DecryptPKCS1v15(rand.Reader, r.priKey, encrypted[i * maxBlockSize : (i + 1) * maxBlockSize])
        if err != nil {
            return nil, err
        }
        buf.Write(tmp)
    }
    return buf.Bytes(), err
}

func (r *Rsa) SignString(src string, base64Encoding bool) (string, error) {
    tmp, err := r.Sign([]byte(src))
    if err != nil {
        return "", err
    }
    return utils.Bytes2String(tmp, base64Encoding), err
}

func (r *Rsa) VerifyString(src, sign string, base64Encoding bool) error {
    tmp, err := utils.String2Bytes(sign, base64Encoding)
    if err != nil {
        return err
    }
    return r.Verify([]byte(src), tmp)
}

func (r *Rsa) Sign(src []byte) ([]byte, error) {
    return r.Sign2(src, crypto.SHA256)
}

func (r *Rsa) Sign2(src []byte, hash crypto.Hash)  ([]byte, error) {
    h := hash.New()
    h.Write(src)
    hashed := h.Sum(nil)
    return rsa.SignPKCS1v15(rand.Reader, r.priKey, hash, hashed)
}

func (r *Rsa) Verify(src, sign []byte) error {
    return r.Verify2(src, sign, crypto.SHA256)
}

func (r *Rsa) Verify2(src, sign []byte, hash crypto.Hash) error {
    h := hash.New()
    h.Write(src)
    hashed := h.Sum(nil)
    return rsa.VerifyPKCS1v15(r.pubKey, hash, hashed, sign)
}

func (r *Rsa) SaveKeyPair(pubKeyFile, priKeyFile string) error {
    err := r.SavePublicKey(pubKeyFile)
    if err != nil {
        return err
    }
    err = r.SavePrivateKey(priKeyFile)
    return err
}

func (r *Rsa) SavePublicKey(pubKeyFile string) error {
    pubKeyBytes, err := x509.MarshalPKIXPublicKey(&(r.priKey).PublicKey)
    if err != nil {
        return err
    }
    block := &pem.Block{Type:"PUBLIC-KEY", Bytes:pubKeyBytes}
    file, err := os.Create(pubKeyFile)
    defer file.Close()
    if err != nil {
        return err
    }
    err = pem.Encode(file, block)
    return err
}

func (r *Rsa) SavePrivateKey(priKeyFile string) error {
    priKeyBytes, err := x509.MarshalPKCS8PrivateKey(r.priKey)
    block := &pem.Block{Type:"PRIVATE-KEY", Bytes:priKeyBytes}
    file, err := os.Create(priKeyFile)
    defer file.Close()
    if err != nil {
        return err
    }
    err = pem.Encode(file, block)
    return err
}

func (r *Rsa) LoadKeyPair(pubKeyFile, priKeyFile string) error {
    err := r.LoadPublicKey(pubKeyFile)
    if err != nil {
        return err
    }
    err = r.LoadPrivateKey(priKeyFile)
    return err
}

func (r *Rsa) LoadPublicKey(pubKeyFile string) error {
    pubKeyBytes, err := ioutil.ReadFile(pubKeyFile)
    if err != nil {
        return err
    }
    block, _ := pem.Decode(pubKeyBytes)
    if block == nil {
        return errors.New(utils.Concat("can't load public key from file \"", pubKeyFile, "\""))
    }
    pub, err := x509.ParsePKIXPublicKey(block.Bytes)
    if err != nil {
        return err
    }
    r.pubKey = pub.(*rsa.PublicKey)
    return err
}

func (r *Rsa) LoadPrivateKey(priKeyFile string) error {
    priKeyBytes, err := ioutil.ReadFile(priKeyFile)
    if err != nil {
        return err
    }
    block, _ := pem.Decode(priKeyBytes)
    if block == nil {
        return errors.New(utils.Concat("can't load private key from file \"", priKeyFile, "\""))
    }
    pri, err := x509.ParsePKCS8PrivateKey(block.Bytes)
    if err != nil {
        return err
    }
    r.priKey = pri.(*rsa.PrivateKey)
    return err
}
