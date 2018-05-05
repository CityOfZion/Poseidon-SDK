package bip32

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha512"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	crypto "multicrypt/crypto"
	elliptic "multicrypt/crypto/elliptic"
)

const (
	FirstHardenedChild        = uint32(0x80000000)
	PublicKeyCompressedLength = 33
)

var (
	PrivateWalletVersion, _ = hex.DecodeString("0488ADE4")
	PublicWalletVersion, _  = hex.DecodeString("0488B21E")
)

// Represents a bip32 extended key containing key data, chain code, parent information, and other meta data
type Key struct {
	Version     []byte // 4 bytes
	Depth       byte   // 1 bytes
	ChildNumber []byte // 4 bytes
	FingerPrint []byte // 4 bytes
	ChainCode   []byte // 32 bytes
	Key         []byte // 33 bytes
	IsPrivate   bool   // unserialized
}

var MasterSecret = "Bitcoin seed"

func init() {
	elliptic.ChosenCurve.SetCurveSecp256r1()
}

// Creates a new master extended key from a seed
func NewMasterKey(seed []byte) (*Key, error) {
	// Generate key and chaincode

	hmac := hmac.New(sha512.New, []byte(MasterSecret))
	hmac.Write([]byte(seed))
	intermediary := hmac.Sum(nil)

	// Split it into our key and chain code
	keyBytes := intermediary[:32]
	chainCode := intermediary[32:]

	// Validate key
	err := validatePrivateKey(keyBytes)
	if err != nil {
		return nil, err
	}

	// Create the key struct
	key := &Key{
		Version:     PrivateWalletVersion,
		ChainCode:   chainCode,
		Key:         keyBytes,
		Depth:       0x0,
		ChildNumber: []byte{0x00, 0x00, 0x00, 0x00},
		FingerPrint: []byte{0x00, 0x00, 0x00, 0x00},
		IsPrivate:   true,
	}

	return key, nil
}

// Derives a child key from a given parent as outlined by bip32
func (key *Key) NewChildKey(childIdx uint32) (*Key, error) {
	hardenedChild := childIdx >= FirstHardenedChild
	childIndexBytes := uint32Bytes(childIdx)

	// Fail early if trying to create hardned child from public key
	if !key.IsPrivate && hardenedChild {
		return nil, errors.New("Can't create hardened child for public key")
	}

	// Get intermediary to create key and chaincode from
	// Hardened children are based on the private key
	// NonHardened children are based on the public key
	var data []byte
	if hardenedChild {
		data = append([]byte{0x0}, key.Key...)
	} else {
		data = publicKeyForPrivateKey(key.Key)
	}
	data = append(data, childIndexBytes...)

	hmac := hmac.New(sha512.New, key.ChainCode)
	hmac.Write(data)
	intermediary := hmac.Sum(nil)

	// Create child Key with data common to all both scenarios
	childKey := &Key{
		ChildNumber: childIndexBytes,
		ChainCode:   intermediary[32:],
		Depth:       key.Depth + 1,
		IsPrivate:   key.IsPrivate,
	}

	// Bip32 CKDpriv
	if key.IsPrivate {
		childKey.Version = PrivateWalletVersion
		hash160Key, _ := crypto.Hash160(publicKeyForPrivateKey(key.Key))
		childKey.FingerPrint = hash160Key[:4]
		childKey.Key = addPrivateKeys(intermediary[:32], key.Key)

		// Validate key
		err := validatePrivateKey(childKey.Key)
		if err != nil {
			return nil, err
		}
		// Bip32 CKDpub
	} else {
		keyBytes := publicKeyForPrivateKey(intermediary[:32])

		// Validate key
		err := validateChildPublicKey(keyBytes)
		if err != nil {
			return nil, err
		}

		childKey.Version = PublicWalletVersion
		hash160Key, _ := crypto.Hash160(key.Key)
		childKey.FingerPrint = hash160Key[:4]
		childKey.Key = addPublicKeys(keyBytes, key.Key)
	}

	return childKey, nil
}

// Create public version of key or return a copy; 'Neuter' function from the bip32 spec
func (key *Key) PublicKey() *Key {
	keyBytes := key.Key

	if key.IsPrivate {
		keyBytes = publicKeyForPrivateKey(keyBytes)
	}

	return &Key{
		Version:     PublicWalletVersion,
		Key:         keyBytes,
		Depth:       key.Depth,
		ChildNumber: key.ChildNumber,
		FingerPrint: key.FingerPrint,
		ChainCode:   key.ChainCode,
		IsPrivate:   false,
	}
}

// Serialized an Key to a 78 byte byte slice
func (key *Key) Serialize() []byte {
	// Private keys should be prepended with a single null byte
	keyBytes := key.Key
	if key.IsPrivate {
		keyBytes = append([]byte{0x0}, keyBytes...)
	}

	// Write fields to buffer in order
	buffer := new(bytes.Buffer)
	buffer.Write(key.Version)
	buffer.WriteByte(key.Depth)
	buffer.Write(key.FingerPrint)
	buffer.Write(key.ChildNumber)
	buffer.Write(key.ChainCode)
	buffer.Write(keyBytes)

	// Append the standard doublesha256 checksum
	serializedKey := addChecksumToBytes(buffer.Bytes())

	return serializedKey
}

// Encode the Key in the standard Bitcoin base58 encoding
func (key *Key) String() string {
	return string(crypto.Base58Encode(key.Serialize()))
}

// Cryptographically secure seed
func NewSeed() ([]byte, error) {
	// Well that easy, just make go read 256 random bytes into a slice
	s := make([]byte, 256)
	_, err := rand.Read([]byte(s))
	return s, err
}

// ECC functions

func addPrivateKeys(key1 []byte, key2 []byte) []byte {
	var key1Int big.Int
	var key2Int big.Int
	key1Int.SetBytes(key1)
	key2Int.SetBytes(key2)

	key1Int.Add(&key1Int, &key2Int)
	key1Int.Mod(&key1Int, elliptic.ChosenCurve.N)

	b := key1Int.Bytes()
	if len(b) < 32 {
		extra := make([]byte, 32-len(b))
		b = append(extra, b...)
	}

	return b
}
func publicKeyForPrivateKey(key []byte) []byte {
	z := new(big.Int)
	z.SetBytes(key)
	pubKey := elliptic.ChosenCurve.ScalarBaseMult(z)
	return compressPublicKey(pubKey.X, pubKey.Y)
}

func addPublicKeys(key1 []byte, key2 []byte) []byte {
	x1, y1 := expandPublicKey(key1)
	x2, y2 := expandPublicKey(key2)

	P := elliptic.Point{x1, y1}
	Q := elliptic.Point{x2, y2}

	result := elliptic.ChosenCurve.Add(P, Q)

	return compressPublicKey(result.X, result.Y)
}

func expandPublicKey(key []byte) (*big.Int, *big.Int) {
	Y := big.NewInt(0)
	X := big.NewInt(0)
	qPlus1Div4 := big.NewInt(0)
	X.SetBytes(key[1:])

	// y^2 = x^3 + ax^2 + b
	// a = 0
	// => y^2 = x^3 + b
	ySquared := X.Exp(X, big.NewInt(3), nil)
	ySquared.Add(ySquared, elliptic.ChosenCurve.B)

	qPlus1Div4.Add(elliptic.ChosenCurve.P, big.NewInt(1))
	qPlus1Div4.Div(qPlus1Div4, big.NewInt(4))

	// sqrt(n) = n^((q+1)/4) if q = 3 mod 4
	Y.Exp(ySquared, qPlus1Div4, elliptic.ChosenCurve.P)

	if uint32(key[0])%2 == 0 {
		Y.Sub(elliptic.ChosenCurve.P, Y)
	}

	return X, Y
}

func uint32Bytes(i uint32) []byte {
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, i)
	return bytes
}

func validatePrivateKey(key []byte) error {
	if fmt.Sprintf("%x", key) == "0000000000000000000000000000000000000000000000000000000000000000" || //if the key is zero
		bytes.Compare(key, elliptic.ChosenCurve.N.Bytes()) >= 0 || //or is outside of the curve
		len(key) != 32 { //or is too short
		return errors.New("Invalid seed")
	}

	return nil
}

func validateChildPublicKey(key []byte) error {
	x, y := expandPublicKey(key)

	if x.Sign() == 0 || y.Sign() == 0 {
		return errors.New("Invalid public key")
	}

	return nil
}

func compressPublicKey(x *big.Int, y *big.Int) []byte {
	var key bytes.Buffer

	// Write header; 0x2 for even y value; 0x3 for odd
	key.WriteByte(byte(0x2) + byte(y.Bit(0)))

	// Write X coord; Pad the key so x is aligned with the LSB. Pad size is key length - header size (1) - xBytes size
	xBytes := x.Bytes()
	for i := 0; i < (PublicKeyCompressedLength - 1 - len(xBytes)); i++ {
		key.WriteByte(0x0)
	}
	key.Write(xBytes)

	return key.Bytes()
}

func addChecksumToBytes(data []byte) []byte {
	checksum, _ := crypto.Checksum(data)
	return append(data, checksum...)
}
