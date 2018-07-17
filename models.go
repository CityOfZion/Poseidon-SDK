package hello

type WrapperContractTransaction struct {
	Inputs []struct {
		External int    `json:"External"`
		PrevHash string `json:"PrevHash"`
		Address  string `json:"Address"`
		Vout     int    `json:"vout"`
		Internal int    `json:"Internal"`
	} `json:"Inputs"`
	Outputs []struct {
		Address string `json:"Address"`
		Value   int64  `json:"Value"`
		AssetID string `json:"AssetID"`
	} `json:"Outputs"`
	Attributes []struct {
		Usage uint16 `json:"Usage"`
		Data  []byte `json:"Data"`
	} `json:"Attributes"`
	Version         int `json:"Version"`
	TransactionType int `json:"TransactionType"`
}
