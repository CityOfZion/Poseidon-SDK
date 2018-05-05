package neo

import (
	"encoding/hex"
	"fmt"
	crypto "multicrypt/crypto"

	bip32 "multicrypt/crypto/bip32"

	bip39 "multicrypt/crypto/bip39"
	bip44 "multicrypt/crypto/bip44"

	"github.com/o3labs/neo-utils/neoutils/btckey"
)

type Coin struct {
}

func init() {
	btckey.ChosenCurve.SetCurveSecp256r1()
	bip32.MasterSecret = "Nist256p1 seed" // This changes for Neo
}
func (c *Coin) GeneratePrivateKey(mnemonic string, external, index int) *bip32.Key {

	seed := bip39.NewSeed(mnemonic, "")

	fmt.Println("SEED: ", hex.EncodeToString(seed))

	masterKey, err := bip32.NewMasterKey(seed)

	if err != nil {

	}
	coin := bip44.TypeAntshares
	account := uint32(0x80000000) // Hardened First child
	chain := uint32(external)
	address := uint32(index)

	fKey, err := bip44.NewKeyFromMasterKey(masterKey, coin, account, chain, address) //purpose is hardcoded

	return fKey
}

func (c *Coin) PubKeyToAddress(pubKey *bip32.Key) string {

	pub_bytes := pubKey.PublicKey().Key

	//No current documentation on what this does, however it is needed for NEO
	pub_bytes = append([]byte{0x21}, pub_bytes...)
	pub_bytes = append(pub_bytes, 0xAC)

	hash160PubKey, _ := crypto.Hash160(pub_bytes)

	versionHash160PubKey := append([]byte{0x17}, hash160PubKey...)

	checksum, _ := crypto.Checksum(versionHash160PubKey)

	checkVersionHash160 := append(versionHash160PubKey, checksum...)

	address := crypto.Base58Encode(checkVersionHash160)

	return address

}

func ToNeoAddress(k bip32.Key) {

	// pubKey := k.PublicKey().Key

}

func ToNeoSignature(pubKey []byte) (signature []byte) {

	pubBytes := append([]byte{0x21}, pubKey...)
	pubBytes = append(pubBytes, 0xAC)

	pubHash, _ := crypto.Hash160(pubBytes)

	return pubHash
}
