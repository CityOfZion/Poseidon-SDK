package signing

import (
	"crypto/sha256"
	"math/big"
	crypto "multicrypt/crypto"
	"multicrypt/crypto/Signing/rfc6979"
	"multicrypt/crypto/bip32"
	elliptic "multicrypt/crypto/elliptic"
)

type PublicKey struct {
	elliptic.Point
}

type PrivateKey struct {
	PublicKey
	D *big.Int
}

// func PrivateKeyFromHexString(key string) ecdsa.PrivateKey {
// 	b, _ := hex.DecodeString(key)

// 	priv := new(ecdsa.PrivateKey)
// 	priv.PublicKey.Curve = elliptic.P256()
// 	priv.D = new(big.Int).SetBytes(b)
// 	priv.PublicKey.X, priv.PublicKey.Y = priv.PublicKey.Curve.ScalarBaseMult(b)
// 	return *priv
// }

func Sign(curve elliptic.EllipticCurve, data []byte, key bip32.Key) ([]byte, error) {

	// We can take the Key as a string or as a byte set. We could just deal with a BIP32 Key and leave it there

	digest, _ := crypto.HashSha256(data)

	r, s, err := rfc6979.SignECDSA(curve, key, digest[:], sha256.New)
	if err != nil {
		return nil, err
	}

	curveOrderByteSize := curve.P.BitLen() / 8
	rBytes, sBytes := r.Bytes(), s.Bytes()
	signature := make([]byte, curveOrderByteSize*2)
	copy(signature[curveOrderByteSize-len(rBytes):], rBytes)
	copy(signature[curveOrderByteSize*2-len(sBytes):], sBytes)

	return signature, nil
}

func Verify(curve elliptic.EllipticCurve, pubKey bip32.Key, signature []byte, hash []byte) bool {
	// TODO: IMPLEMENT THIS FEATURE
	return true
}
