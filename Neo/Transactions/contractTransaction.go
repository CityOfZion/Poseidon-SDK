package transactions

type Input struct {
	prevHash  string
	prevIndex uint8
}
type Output struct {
	assetId string
	value   float32
	address string
}

type ContractTransaction struct {
	Id         int
	Attributes []string
	Inputs     []Input
	Outputs    []Output
	Scripts    []string
	SystemFee  float32
	NetworkFee float32
}

func CreateContractTransaction() *ContractTransaction {
	tx := ContractTransaction{
		Id:         128,
		Attributes: []string{},
		Inputs:     []Input{},
		Outputs:    []Output{},
		Scripts:    []string{},
		SystemFee:  0.0,
		NetworkFee: 0.0,
	}

	return &tx
}

func (c *ContractTransaction) addAttributes() {
	// Things like remarks will go here
}

func (c *ContractTransaction) AddOutput(assetId string, value float32, address string) {
	c.Outputs = append(c.Outputs, Output{assetId, value, address})
}

func (c *ContractTransaction) Hash() string {
	return "Hash of the transaction as it is the ID"
}

func (c *ContractTransaction) SetSystemFee(value float32) {
	c.SystemFee = 0.0
	// Contract transaction fees are always zero
}

func (c *ContractTransaction) SetNetworkFee(value float32) {
	c.NetworkFee = value
}

func (c *ContractTransaction) GetUTXOs(seed string) {
	// Method will use seed to get all accounts that have funds

	// Then get together enough inputs to cover the outputs provided
}

func (c *ContractTransaction) signTransaction(key string) {
	// get the contract signing method to sign the transaction

}

func SerialiseTransaction() string {
	return "Serialised transaction"
}
