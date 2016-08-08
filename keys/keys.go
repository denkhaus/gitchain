package keys

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/gob"
	"math/big"

	"code.google.com/p/go.crypto/ripemd160"

	"github.com/conformal/btcec"
	"github.com/denkhaus/gitchain/util"
	"github.com/tv42/base58"
)

// For now, ECDSA keys generated by Gitchain use the P-256 curve
// There are different opinions about what curves to use:
//
//  http://safecurves.cr.yp.to/
//  http://infosecurity.ch/20100926/not-every-elliptic-curve-is-the-same-trough-on-ecc-security/
//  http://www.hyperelliptic.org/tanja/vortraege/20130531.pdf
//
// For now I decided to stick to one used in Bitcoin (secp256k1)
func GenerateECDSA() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(btcec.S256(), rand.Reader)
}

func EncodeECDSAPrivateKey(key *ecdsa.PrivateKey) ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode([]big.Int{*key.X, *key.Y, *key.D})
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func DecodeECDSAPrivateKey(b []byte) (*ecdsa.PrivateKey, error) {
	var p []big.Int
	buf := bytes.NewBuffer(b)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&p)
	if err != nil {
		return nil, err
	}
	privateKey := new(ecdsa.PrivateKey)
	privateKey.PublicKey.Curve = btcec.S256()
	privateKey.PublicKey.X = &p[0]
	privateKey.PublicKey.Y = &p[1]
	privateKey.D = &p[2]
	return privateKey, nil
}

func ECDSAPublicKeyToString(key ecdsa.PublicKey) string {
	x := key.X.Bytes()
	y := key.Y.Bytes()
	sha := util.SHA256(append(append([]byte{0x04}, x...), y...)) // should it be 0x04?
	ripe := ripemd160.New().Sum(sha)
	ripesha := util.SHA256(ripe)
	ripedoublesha := util.SHA256(ripesha)
	head := ripedoublesha[0:3]
	final := append(ripe, head...)
	i := new(big.Int)
	i.SetBytes(final)
	var b []byte
	return string(base58.EncodeBig(b, i))
}

func EncodeECDSAPublicKey(key *ecdsa.PublicKey) ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode([]big.Int{*key.X, *key.Y})
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func DecodeECDSAPublicKey(b []byte) (*ecdsa.PublicKey, error) {
	var p []big.Int
	buf := bytes.NewBuffer(b)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&p)
	if err != nil {
		return nil, err
	}
	return &ecdsa.PublicKey{Curve: btcec.S256(), X: &p[0], Y: &p[1]}, nil
}

func EqualECDSAPrivateKeys(k1, k2 *ecdsa.PrivateKey) (bool, error) {
	k1e, err := EncodeECDSAPrivateKey(k1)
	if err != nil {
		return false, err
	}
	k2e, err := EncodeECDSAPrivateKey(k2)
	if err != nil {
		return false, err
	}
	return bytes.Compare(k1e, k2e) == 0, nil
}
