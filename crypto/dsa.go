package crypto

import (
    "crypto/dsa"
    "crypto/rand"
    "crypto"
    "math/big"
    "github.com/blazecrystal/beyondts-go/utils"
)

type DSA struct {
    priKey *dsa.PrivateKey
    pubKey *dsa.PublicKey
}

func NewDSAInstance() *DSA {
    return &DSA{}
}

func (d *DSA) GenKeyPair(paramSize dsa.ParameterSizes) error {
    d.priKey = &dsa.PrivateKey{}
    err := dsa.GenerateParameters(&d.priKey.Parameters, rand.Reader, paramSize)
    if err != nil {
        return err
    }
    err = dsa.GenerateKey(d.priKey, rand.Reader)
    d.pubKey = &d.priKey.PublicKey
    return err;
}

func (d *DSA) SignString(src string, base64Encoding bool) (r, s string, err error) {
    tmpR, tmpS, err := d.Sign([]byte(src))
    if err != nil {
        return "", "", err
    }
    return utils.Bytes2String(tmpR.Bytes(), base64Encoding), utils.Bytes2String(tmpS.Bytes(), base64Encoding), err
}

func (d *DSA) Sign(src []byte) (r, s *big.Int, err error) {
    return d.Sign2(src, crypto.SHA256)
}

func (d *DSA) Sign2(src []byte, hash crypto.Hash) (r, s *big.Int, err error) {
    h := hash.New()
    h.Write(src)
    hashed := h.Sum(nil)
    return dsa.Sign(rand.Reader, d.priKey, hashed)
}

func (d *DSA) VerifyString(src, r, s string, base64Encoding bool) (bool, error) {
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
    return d.Verify([]byte(src), ri, si), err
}

func (d *DSA) Verify(src []byte, r, s *big.Int) bool {
    return d.Verify2(src, crypto.SHA256, r, s)
}

func (d *DSA) Verify2(src []byte, hash crypto.Hash, r, s *big.Int) bool {
    h := hash.New()
    h.Write(src)
    hashed := h.Sum(nil)
    return dsa.Verify(d.pubKey, hashed, r, s)
}
