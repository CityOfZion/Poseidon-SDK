package transactions

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateContractTransaction(t *testing.T) {
	tx1 := CreateContractTransaction()
	tx1.AddOutput("NEO", 12, "address1")
	tx1.AddOutput("GAS", 1, "Address2")
	tx1.SetSystemFee()
	fmt.Println(tx1.SystemFee)
	t.Fail()
	assert.Equal(t, "KxcqV28rGDcpVR3fYg7R9vricLpyZ8oZhopyFLAWuRv7Y8TE9WhW", "crypto.WIFEncode(privKey.Key)")
}
