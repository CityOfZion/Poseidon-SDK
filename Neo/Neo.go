package neo

import (
	crypto "multicrypt/crypto"

	bip32 "multicrypt/crypto/bip32"

	bip39 "multicrypt/crypto/bip39"
	bip44 "multicrypt/crypto/bip44"

	"github.com/o3labs/neo-utils/neoutils/btckey"
)

type Coin struct{}

func init() {
	btckey.ChosenCurve.SetCurveSecp256r1()
	bip32.MasterSecret = "Nist256p1 seed" // This changes for Neo
}
func (c *Coin) GeneratePrivateKey(mnemonic string, external, index int) *bip32.Key {

	seed := bip39.NewSeed(mnemonic, "")

	masterKey, err := bip32.NewMasterKey(seed)

	if err != nil {
		return nil
	}
	coin := bip44.TypeAntshares
	account := uint32(0x80000000) // Hardened First child
	chain := uint32(external)
	address := uint32(index)

	fKey, err := bip44.NewKeyFromMasterKey(masterKey, coin, account, chain, address) //purpose is hardcoded

	return fKey
}

func (c *Coin) PubKeyToAddress(pubKey *bip32.Key) string {

	publicKeyBytes := pubKey.PublicKey().Key

	//No current documentation on what this does, however it is needed for NEO
	// This is for the scriptHash, the verification script
	// This could be added at the time of creating the verif script, so no reason for this to be here
	// The 0xAC is CHECKSIG
	publicKeyBytes = append([]byte{0x21}, publicKeyBytes...)
	publicKeyBytes = append(publicKeyBytes, 0xAC)

	hash160PubKey, _ := crypto.Hash160(publicKeyBytes)

	versionHash160PubKey := append([]byte{0x17}, hash160PubKey...)

	checksum, _ := crypto.Checksum(versionHash160PubKey)

	checkVersionHash160 := append(versionHash160PubKey, checksum...)

	address := crypto.Base58Encode(checkVersionHash160)

	return address

}

func (c *Coin) GenerateMnemonic() string {
	randomNum, _ := bip39.NewEntropy(128)
	mnemonic, _ := bip39.NewMnemonic(randomNum)
	return mnemonic
}
