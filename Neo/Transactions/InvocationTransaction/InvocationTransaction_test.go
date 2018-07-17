package InvocationTransaction

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

func TestTransfer(t *testing.T) {

	tx := CreateInvocationTransaction()
	tx.AppendScriptAttribute("ALxUuc4gSNaYdUUrPudem4DMkowz6x6Rwo")
	// tx.Script.SetOperation("transfer")
	// tx.Script.AddArgument(tx.AToS("ALxUuc4gSNaYdUUrPudem4DMkowz6x6Rwo"))
	tx.Script.SetScriptHash("9aff1e08aea2048a26a3d2ddbb3df495b932b1e7")

	tx.TransferTokens("ALxUuc4gSNaYdUUrPudem4DMkowz6x6Rwo", "AQGV3Q4Q4kczcb7Q6Ax11zBP5WcygRRKXw", 100)
	// tx.AddAttributes(uint16(0xff), []byte(time.Now().String())) //
	tx.AddPrivateKey("online ramp onion faculty trap clerk near rabbit busy gravity prize employ", 0, 1)
	tx.SignTransaction("online ramp onion faculty trap clerk near rabbit busy gravity prize employ")
	buf := new(bytes.Buffer)
	err := tx.SerialiseTransaction(buf, true)

	if err != nil {
		t.Fail()
		fmt.Println("Err", err)
	}
	expectedRawTrans := "d1014c0164145d2a330652006323d660fe74a28ca4a46a851ef21438da38f5350186eff4f0d1f086832aec77b8f8cd53c1087472616e7366657267e7b132b995f43dbbddd2a3268a04a2ae081eff9a0000000000000000012038da38f5350186eff4f0d1f086832aec77b8f8cd0000014140ecaf3cbd2438f5e84d5f48b1131b4f85d2e2661905550c016eb1a999df2d2d47ee96b7b70d9a89867ea9c4c98a89a820c8bef5fa79c61efd42a9554fd75da773232102ada2c1dd51b932dde3218181cf53751adee33d8731abb05f7e44d9aff2fb1026ac"
	expectedTxID := "7e0c4bb50d49402a73af03b5b5bf06d1378ea5db88aa0f41dfd8e3f42ec4852d"
	assert.Equal(t, expectedRawTrans, hex.EncodeToString(buf.Bytes()))
	assert.Equal(t, expectedTxID, tx.GetHash())
}
