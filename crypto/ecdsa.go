package crypto

import (
    "crypto/rand"
    "crypto"
    "math/big"
    "github.com/blazecrystal/beyondts-go/utils"
    "errors"
    "crypto/x509"
    "io/ioutil"
    "encoding/pem"
    "os"
    "crypto/ecdsa"
    "crypto/elliptic"
)

type ECDSA struct {
    priKey *ecdsa.PrivateKey
    pubKey *ecdsa.PublicKey
}

func NewECDSAInstance() *ECDSA {
    return &ECDSA{}
}

func (e *ECDSA) GenKeyPair(c elliptic.Curve) error {
    tmp, err := ecdsa.GenerateKey(c, rand.Reader)
    if err != nil {
        return err
    }
    e.priKey = tmp
    e.pubKey = &e.priKey.PublicKey
    return err;
}

func (e *ECDSA) SignString(src string, base64Encoding bool) (r, s string, err error) {
    tmpR, tmpS, err := e.Sign([]byte(src))
    if err != nil {
        return "", "", err
    }
    return utils.Bytes2String(tmpR.Bytes(), base64Encoding), utils.Bytes2String(tmpS.Bytes(), base64Encoding), err
}

func (e *ECDSA) Sign(src []byte) (r, s *big.Int, err error) {
    return e.Sign2(src, crypto.SHA256)
}

func (e *ECDSA) Sign2(src []byte, hash crypto.Hash) (r, s *big.Int, err error) {
    h := hash.New()
    h.Write(src)
    hashed := h.Sum(nil)
    return ecdsa.Sign(rand.Reader, e.priKey, hashed)
}

func (e *ECDSA) VerifyString(src, r, s string, base64Encoding bool) (bool, error) {
    tmpR, err := utils.String2Bytes(r, base64Encoding)
    if err != nil {
        return false, err
    }
    tmpS, err := utils.String2Bytes(s, base64Encoding)
    if err != nil {
        return false, err
    }
    ri := new(big.Int).SetBytes(tmpR)
    si := new(big.Int).SetBytes(tmpS)
    return e.Verify([]byte(src), ri, si), err
}

func (e *ECDSA) Verify(src []byte, r, s *big.Int) bool {
    return e.Verify2(src, crypto.SHA256, r, s)
}

func (e *ECDSA) Verify2(src []byte, hash crypto.Hash, r, s *big.Int) bool {
    h := hash.New()
    h.Write(src)
    hashed := h.Sum(nil)
    return ecdsa.Verify(e.pubKey, hashed, r, s)
}

func (e *ECDSA) SaveKeyPair(pubKeyFile, priKeyFile string) error {
    err := e.SavePublicKey(pubKeyFile)
    if err != nil {
        return err
    }
    err = e.SavePrivateKey(priKeyFile)
    return err
}

func (e *ECDSA) SavePublicKey(pubKeyFile string) error {
    pubKeyBytes, err := x509.MarshalPKIXPublicKey(&(e.priKey).PublicKey)
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

func (e *ECDSA) SavePrivateKey(priKeyFile string) error {
    priKeyBytes, err := x509.MarshalPKCS8PrivateKey(e.priKey)
    block := &pem.Block{Type:"PRIVATE-KEY", Bytes:priKeyBytes}
    file, err := os.Create(priKeyFile)
    defer file.Close()
    if err != nil {
        return err
    }
    err = pem.Encode(file, block)
    return err
}

func (e *ECDSA) LoadKeyPair(pubKeyFile, priKeyFile string) error {
    err := e.LoadPublicKey(pubKeyFile)
    if err != nil {
        return err
    }
    err = e.LoadPrivateKey(priKeyFile)
    return err
}

func (e *ECDSA) LoadPublicKey(pubKeyFile string) error {
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
    e.pubKey = pub.(*ecdsa.PublicKey)
    return err
}

func (e *ECDSA) LoadPrivateKey(priKeyFile string) error {
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
    e.priKey = pri.(*ecdsa.PrivateKey)
    return err
}