package btc

import (
	"encoding/hex"
	"fmt"
	crypto "multicrypt/crypto"

	bip32 "github.com/FactomProject/go-bip32"
	bip39 "github.com/FactomProject/go-bip39"
	bip44 "github.com/FactomProject/go-bip44"
	"github.com/o3labs/neo-utils/neoutils/btckey"
)

type Coin struct{}

func init() {
	btckey.ChosenCurve.SetCurveSecp256k1()
}
func (c *Coin) GeneratePrivateKey(mnemonic string, external, index int) string {

	seed := bip39.NewSeed(mnemonic, "")

	masterKey, err := bip32.NewMasterKey(seed)

	if err != nil {

	}
	coin := uint32(0x80000000)    // Coin number for NEO
	account := uint32(0x80000000) // Hardened First child
	chain := uint32(external)
	address := uint32(index)
	fKey, err := bip44.NewKeyFromMasterKey(masterKey, coin, account, chain, address) //purpose is hardcoded

	return hex.EncodeToString(fKey.Key)
}

func (c *Coin) PrivateKeyToAddress(privKey string) string {
	pb, _ := hex.DecodeString(privKey)
	var priv btckey.PrivateKey
	err := priv.FromBytes(pb) // convert private key to public key
	if err != nil {
		return "Error"
	}
	pub_bytes := priv.PublicKey.ToBytes()

	fmt.Println("PUBLIC KEY: ", hex.EncodeToString(pub_bytes))

	// fmt.Println(hex.EncodeToString(priv.PublicKey.ToBytes()))
	hash160PubKey, _ := crypto.Hash160(pub_bytes)

	versionHash160PubKey := append([]byte{0x00}, hash160PubKey...)

	checksum, _ := crypto.Checksum(versionHash160PubKey)

	checkVersionHash160 := append(versionHash160PubKey, checksum...)

	address := crypto.Base58Encode(checkVersionHash160)

	return address
	// return ""
}
