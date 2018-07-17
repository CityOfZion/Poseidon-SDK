package ClaimTransaction

import (
	"bytes"
	"encoding/hex"
	"fmt"
	fixed8 "multicrypt/Utils/Fixed8"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	NEO = "c56f33fc6ecfcd0c225c4ab356fee59390af8560be0e930faebe74a6daff7c9b"
	GAS = "602c79718b16e442de58778e148d0b1084e3b2dffd5de6b7b16cee7969282de7"
)

func TestCreateClaimTransaction(t *testing.T) {
	tx := CreateClaimTransaction()
	tx.AddClaim("480a72a7d43452f9503f4f9552768cfd54ffa103776b703e7519d22d959ade68", 0, "ALxUuc4gSNaYdUUrPudem4DMkowz6x6Rwo", 0, 1)
	val, _ := fixed8.Fixed8DecodeString("0.02119752")
	tx.AddOutput(GAS, int64(val), "ALxUuc4gSNaYdUUrPudem4DMkowz6x6Rwo")
	tx.AddPrivateKey("online ramp onion faculty trap clerk near rabbit busy gravity prize employ", 0, 1)
	tx.SignTransaction("online ramp onion faculty trap clerk near rabbit busy gravity prize employ")
	buf := new(bytes.Buffer)
	err := tx.SerialiseTransaction(buf, true)
	if err != nil {
		fmt.Println("Err", err)
	}

	expected := "02000168de9a952dd219753e706b7703a1ff54fd8c7652954f3f50f95234d4a7720a480000000001e72d286979ee6cb1b7e65dfddfb2e384100b8d148e7758de42e4168b71792c60485820000000000038da38f5350186eff4f0d1f086832aec77b8f8cd01414032fa76069ef7b717b8f1ec02883f915674a8932f1a0ca279cf3ecb6a20e5573d761058e9f8eb97d120efd498f1858e4c51533eb04415a83ab902483b5c9b4f3a232102ada2c1dd51b932dde3218181cf53751adee33d8731abb05f7e44d9aff2fb1026ac"
	assert.Equal(t, expected, hex.EncodeToString(buf.Bytes()))

}
