package ContractTransaction

import (
	"io"
	base "multicrypt/Neo/Transactions"
)

type ContractTransaction struct {
	base.BasicTransaction
}

func CreateContractTransaction() *ContractTransaction {
	tx := ContractTransaction{
		base.BasicTransaction{
			Type:       128,
			Version:    0,
			Attributes: []base.Attributes{},
			Inputs:     []base.Input{},
			Outputs:    []base.Output{},
			Witnesses:  []base.Witness{},
			SystemFee:  0.0,
			NetworkFee: 0.0,
		},
	}
	tx.F = tx.serialiseAdditionalFields
	return &tx
}

func (c *ContractTransaction) serialiseAdditionalFields(w io.Writer) error {
	return nil
}
