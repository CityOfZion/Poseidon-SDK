package ClaimTransaction

import (
	"io"
	base "multicrypt/Neo/Transactions"
	writer "multicrypt/Utils/Writer"
)

type ClaimTransaction struct {
	base.BasicTransaction
	Claims []base.Input
}

func CreateClaimTransaction() *ClaimTransaction {
	tx := ClaimTransaction{
		base.BasicTransaction{
			Type:       2,
			Version:    0,
			Attributes: []base.Attributes{},
			Inputs:     []base.Input{},
			Outputs:    []base.Output{},
			Witnesses:  []base.Witness{},
			SystemFee:  0.0,
			NetworkFee: 0.0,
		},
		[]base.Input{},
	}
	tx.F = tx.serialiseAdditionalFields
	return &tx
}
func (c *ClaimTransaction) serialiseAdditionalFields(w io.Writer) error {
	if err := writer.WriteVarUint(w, uint64(len(c.Claims))); err != nil {
		return err
	}
	for _, input := range c.Claims {
		if err := input.Encode(w); err != nil {
			return err
		}
	}
	return nil
}
func (c *ClaimTransaction) AddClaim(prevHash string, prevIndex uint16, Address string, External, Index int) {
	c.Claims = append(c.Claims, base.Input{prevHash, prevIndex, Address, External, Index})

}
