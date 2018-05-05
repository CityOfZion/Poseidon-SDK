package neo

import (
	"encoding/hex"
	network "multicrypt/Neo/API"
	"multicrypt/crypto"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWIFEncoder(t *testing.T) {
	// Test taken from : https://coranos.github.io/neo/ledger-nano-s/recovery/
	neoCoin := Coin{}
	mnemonic := "scare daughter hazard climb layer card useful find giraffe play street bonus depend execute appear never book file shock nest strike impulse clarify vintage"
	privKey := neoCoin.GeneratePrivateKey(mnemonic, 0, 0)

	assert.Equal(t, "KxcqV28rGDcpVR3fYg7R9vricLpyZ8oZhopyFLAWuRv7Y8TE9WhW", crypto.WIFEncode(privKey.Key))
}

func TestPrvateKeys(t *testing.T) {
	neoCoin := Coin{}
	mnemonic := "online ramp onion faculty trap clerk near rabbit busy gravity prize employ"
	privKey := neoCoin.GeneratePrivateKey(mnemonic, 0, 0) // returns a bip32.Key

	assert.Equal(t, "a3420cdf2b2e8c1a0612f4db006619d2145234eda84a80e90fc499cb582ded7f", hex.EncodeToString(privKey.Key))
}

func TestPrivateKeyToAddress(t *testing.T) {
	neoCoin := Coin{}
	mnemonic := "online ramp onion faculty trap clerk near rabbit busy gravity prize employ"
	prKey := neoCoin.GeneratePrivateKey(mnemonic, 0, 0)
	address := neoCoin.PubKeyToAddress(prKey.PublicKey())
	assert.Equal(t, "AHDKs3tXxhws9twq2Bk7VL1GGqfmSjEVhV", address)
}

func TestGetBalanceOfUsedAddress(t *testing.T) {
	nw := network.API{}
	response := nw.CheckBalance("AXym3Qc9mRbKF5HWAtmJUoHvB3F8brEdLe")
	assert.NotEqual(t, "not found", response.Address) // Test will fail, if the address has not been used yet
}

func TestGetBalanceOfNewAddress(t *testing.T) {
	nw := network.API{}
	response := nw.CheckBalance("AT9DbViHVMy52FZnw2XDjT3TMkMNYyZZmm")
	assert.Equal(t, "not found", response.Address) // Test will fail, if the address has not been used yet
}
