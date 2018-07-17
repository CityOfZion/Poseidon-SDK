package ContractTransaction

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	NEO = "c56f33fc6ecfcd0c225c4ab356fee59390af8560be0e930faebe74a6daff7c9b"
	GAS = "602c79718b16e442de58778e148d0b1084e3b2dffd5de6b7b16cee7969282de7"
)

func TestCreateContractTransaction(t *testing.T) {

	tx := CreateContractTransaction()

	tx.AddOutput(GAS, int64(100000000000), "ANmQW9femWe1oPgEDD4pSTaNUu5rvHYW1R")
	tx.AddOutput(GAS, int64(100000000000), "ALxUuc4gSNaYdUUrPudem4DMkowz6x6Rwo")
	tx.AddInput("f08dfee4598e3a91e20b151ad63a967b95240c5ec96ecf4548e27dcff7a330d7", 0, "ALxUuc4gSNaYdUUrPudem4DMkowz6x6Rwo", 0, 1)

	tx.SignTransaction("online ramp onion faculty trap clerk near rabbit busy gravity prize employ")

	buf := new(bytes.Buffer)

	err := tx.SerialiseTransaction(buf, true)
	if err != nil {
		fmt.Println("Err", err)
	}

	expectedRawTrans := "80000001d730a3f7cf7de24845cf6ec95e0c24957b963ad61a150be2913a8e59e4fe8df0000002e72d286979ee6cb1b7e65dfddfb2e384100b8d148e7758de42e4168b71792c6000e87648170000004cb238cca4811a0b41cd59b38db35d5d71ad560ee72d286979ee6cb1b7e65dfddfb2e384100b8d148e7758de42e4168b71792c6000e876481700000038da38f5350186eff4f0d1f086832aec77b8f8cd0141407ac3f01b0feaec7afe18b84654dbc0b37473cca3478e8e4d27711f788666fb4965f534d66c13fabcccf936fc1b117f2c1edc20b7658f3d649a00bd698c214412232102ada2c1dd51b932dde3218181cf53751adee33d8731abb05f7e44d9aff2fb1026ac"
	expectedTXID := "9545eade0ea602e349514749e7740de47abdd9b076f333c45e62adb41db0a3f8"

	assert.Equal(t, expectedRawTrans, hex.EncodeToString(buf.Bytes()))
	assert.Equal(t, expectedTXID, tx.GetHash())
}
