package InvocationTransaction

// This package will form the invocation transaction and use rpc to send it off
// We need a scriptbuilder to build the script.

import (
	"encoding/hex"
	"fmt"
	"io"
	scriptBuilder "multicrypt/NEO/Transactions/scriptBuilder"
	base "multicrypt/Neo/Transactions"
	uint160 "multicrypt/Utils/Uint160"
)

type InvocationTransaction struct {
	base.BasicTransaction
	Script scriptBuilder.Script
}

func CreateInvocationTransaction() *InvocationTransaction {
	tx := InvocationTransaction{
		base.BasicTransaction{
			Type:       209,
			Version:    1,
			Attributes: []base.Attributes{},
			Inputs:     []base.Input{},
			Outputs:    []base.Output{},
			Witnesses:  []base.Witness{},
			SystemFee:  0.0,
			NetworkFee: 0.0,
		}, scriptBuilder.Script{},
	}
	tx.F = tx.SerialiseSpecialFields

	return &tx
}

func (i *InvocationTransaction) AppendScriptAttribute(address string) {
	scriptHash := base.AddressToScriptHash(address)

	scriptHashAsBytes, err := hex.DecodeString(scriptHash)
	if err != nil {

		fmt.Println("Err with Input Scripthash conversion")
		return
	}

	i.Attributes = append(i.Attributes, base.Attributes{uint16(0x20), scriptHashAsBytes})
}

//TODO: go into the script build and here and check for errors when writing to binary
func (c *InvocationTransaction) SerialiseSpecialFields(w io.Writer) error {

	c.Script.SerialiseScript(w)
	// MARK: This does not work for more than one script, will look into it

	// lengthOfScripts := uint64(0)
	// for _, script := range c.Scripts {
	// 	lengthOfScripts += script.GetLength()
	// }
	// writer.WriteVarUint(w, uint64(lengthOfScripts))

	// for _, script := range c.Scripts {
	// 	script.SerialiseScript(w)
	// }
	return nil
}

//NEP5
// need an encoding scheme for strings, to verify whether it is an address or it is normal
func (c *InvocationTransaction) TransferTokens(from, to string, amount int64) {

	toScriptHash, _ := uint160.Uint160DecodeString(base.AddressToScriptHash(to))
	fromScriptHash, _ := uint160.Uint160DecodeString(base.AddressToScriptHash(from))

	c.Script.SetOperation("transfer")
	c.Script.AddArgument(amount)
	c.Script.AddArgument(toScriptHash.Bytes())
	c.Script.AddArgument(fromScriptHash.Bytes())

}
func (c *InvocationTransaction) AToS(s string) []byte {
	sh, _ := uint160.Uint160DecodeString(base.AddressToScriptHash(s))
	return sh.Bytes()
}

// TODO: Add other NEP Methods
