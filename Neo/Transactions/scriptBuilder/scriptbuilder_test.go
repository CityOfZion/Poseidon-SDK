package scriptBuilder

import (
	"bytes"
	"encoding/hex"
	"fmt"
	transactions "multicrypt/Neo/Transactions"
	util "multicrypt/Utils/Fixed8"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmitString(t *testing.T) {
	scriptBuilder := Script{}
	b := new(bytes.Buffer)
	scriptBuilder.EmitString(b, "name")

	assert.Equal(t, "046e616d65", hex.EncodeToString(b.Bytes()))
}

func TestEmitArr(t *testing.T) {
	b := new(bytes.Buffer)
	scriptBuilder := Script{}
	arr := []interface{}{}
	scriptBuilder.EmitArr(b, arr)
	assert.Equal(t, "00c1", hex.EncodeToString(b.Bytes()))

	b = new(bytes.Buffer)

	arr = []interface{}{"name"}
	scriptBuilder.EmitArr(b, arr)
	assert.Equal(t, "046e616d6551c1", hex.EncodeToString(b.Bytes()))
}

func TestEmitAppCall(t *testing.T) {
	b := new(bytes.Buffer)
	scriptBuilder := Script{}
	scriptBuilder.EmitAppCall(b)
	assert.Equal(t, "67", hex.EncodeToString(b.Bytes()))
}

func TestEmitBool(t *testing.T) {
	b := new(bytes.Buffer)
	scriptBuilder := Script{}
	scriptBuilder.EmitBool(b, true)
	assert.Equal(t, "51", hex.EncodeToString(b.Bytes()))

	b = new(bytes.Buffer)
	scriptBuilder.EmitBool(b, false)
	assert.Equal(t, "00", hex.EncodeToString(b.Bytes()))
}

func TestEmitInt(t *testing.T) {
	scriptBuilder := Script{}
	b := new(bytes.Buffer)
	scriptBuilder.EmitIntLiteral(b, -1)
	assert.Equal(t, "4f", hex.EncodeToString(b.Bytes()))

	b = new(bytes.Buffer)
	scriptBuilder.EmitIntLiteral(b, 10)
	fmt.Println(hex.EncodeToString(b.Bytes()))
	assert.Equal(t, hex.EncodeToString([]byte{byte(PUSH10)}), hex.EncodeToString(b.Bytes()))

	b = new(bytes.Buffer)
	scriptBuilder.EmitIntLiteral(b, 100)
	assert.Equal(t, "0164", hex.EncodeToString(b.Bytes()))

	b = new(bytes.Buffer)
	scriptBuilder.EmitInt(b, 200100000000)
	assert.Equal(t, "0500b1e3962e", hex.EncodeToString(b.Bytes()))
}

func TestScriptBuilder(t *testing.T) {

	s := Script{
		Operation:  "name",
		Args:       []interface{}{},
		ScriptHash: "5b7074e873973a6ed3708862f219a6fbf4d1c411",
	}
	assert.Equal(t, "00c1046e616d656711c4d1f4fba619f2628870d36e3a9773e874705b", s.Hex())
}

func TestTransfer(t *testing.T) {

	scriptHash1 := transactions.AddressToScriptHash("APmxCmvGgb7bayoAApvFC6Aqa9pNqnUB5N")
	from, _ := hex.DecodeString(scriptHash1)
	scriptHash2 := transactions.AddressToScriptHash("APmxCmvGgb7bayoAApvFC6Aqa9pNqnUB5N")
	to, _ := hex.DecodeString(scriptHash2)
	amount := util.Fixed8(100)
	s := Script{
		Operation:  "transfer",
		Args:       []interface{}{from, to, amount},
		ScriptHash: "9aff1e08aea2048a26a3d2ddbb3df495b932b1e7",
	}
	assert.Equal(t, "1457c4cf51f12ce6d78a585f1ea9bd1f3927c7232c1457c4cf51f12ce6d78a585f1ea9bd1f3927c7232c53c1087472616e7366657267e7b132b995f43dbbddd2a3268a04a2ae081eff9a", s.Hex())

}
