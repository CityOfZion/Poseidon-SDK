package neo

import (
	"encoding/hex"
	"multicrypt/crypto"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWIFEncoder(t *testing.T) {
	// Test taken from : https://coranos.github.io/neo/ledger-nano-s/recovery/
	neoCoin := Coin{}
	mnemonic := "online ramp onion faculty trap clerk near rabbit busy gravity prize employ"
	privKey := neoCoin.GeneratePrivateKey(mnemonic, 0, 1)
	assert.Equal(t, "KyBJR4835UvmiTBaUwvAN7WVdtpMMCkN3Wqnn7F8UFvd6UaVXFuQ", crypto.WIFEncode(privKey.Key))
}

func TestPrivateKeyGeneration(t *testing.T) {
	neoCoin := Coin{}
	mnemonic := "online ramp onion faculty trap clerk near rabbit busy gravity prize employ"
	privKey := neoCoin.GeneratePrivateKey(mnemonic, 0, 0) // returns a bip32.Key

	assert.Equal(t, "a3420cdf2b2e8c1a0612f4db006619d2145234eda84a80e90fc499cb582ded7f", hex.EncodeToString(privKey.Key))
}

func TestPrivateKeyToAddress(t *testing.T) {
	neoCoin := Coin{}
	mnemonic := "online ramp onion faculty trap clerk near rabbit busy gravity prize employ"

	prKey := neoCoin.GeneratePrivateKey(mnemonic, 0, 1)
	address := neoCoin.PubKeyToAddress(prKey.PublicKey())
	assert.Equal(t, "ALxUuc4gSNaYdUUrPudem4DMkowz6x6Rwo", address)
}
func TestPubToAddress(t *testing.T) {

	neoCoin := Coin{}
	mnemonic := "online ramp onion faculty trap clerk near rabbit busy gravity prize employ"
	privKey := neoCoin.GeneratePrivateKey(mnemonic, 0, 1) // returns a bip32.Key
	pubKey := privKey.PublicKey()
	address := neoCoin.PubKeyToAddress(pubKey)
	assert.Equal(t, "ALxUuc4gSNaYdUUrPudem4DMkowz6x6Rwo", address)
}
