package btc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test was done by running the TestWIfi encoder to give back the WIF, which was
// Then put into neon wallet to find the address.
// The address was then used in the second test

func TestPrivateKeyGen(t *testing.T) {
	//https://iancoleman.io/bip39/#english
	btcCoin := Coin{}
	mnemonic := "scare daughter hazard climb layer card useful find giraffe play street bonus depend execute appear never book file shock nest strike impulse clarify vintage"
	address := btcCoin.PrivateKeyToAddress(btcCoin.GeneratePrivateKey(mnemonic, 0, 107))

	assert.Equal(t, "195pRLBnd51BzuiMqVQ4W4SUoQH6fsbrEe", address)
}
