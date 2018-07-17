package rpc

// Models generated from: https://mholt.github.io/json-to-go/

type BlockRes struct {
	Hash              string `json:"hash"`
	Size              int    `json:"size"`
	Version           int    `json:"version"`
	Previousblockhash string `json:"previousblockhash"`
	Merkleroot        string `json:"merkleroot"`
	Time              int    `json:"time"`
	Index             int    `json:"index"`
	Nonce             string `json:"nonce"`
	Nextconsensus     string `json:"nextconsensus"`
	Script            struct {
		Invocation   string `json:"invocation"`
		Verification string `json:"verification"`
	} `json:"script"`
	Tx []struct {
		Txid       string        `json:"txid"`
		Size       int           `json:"size"`
		Type       string        `json:"type"`
		Version    int           `json:"version"`
		Attributes []interface{} `json:"attributes"`
		Vin        []struct {
			Txid      string `json:"txid"`
			VoutIndex int    `json:"vout"`
		} `json:"vin"`
		Vout []struct {
			N       int    `json:"n"`
			Address string `json:"address"`
			Value   string `json:"value"`
			Asset   string `json:"asset"`
		} `json:"vout"`
		SysFee  string `json:"sys_fee"`
		NetFee  string `json:"net_fee"`
		Scripts []struct {
			VerificationScript string `json:"verification"`
			InvocationScript   string `json:"invocation"`
		} `json:"scripts"`
		Nonce int64 `json:"nonce"`
	} `json:"tx"`
	Confirmations int    `json:"confirmations"`
	Nextblockhash string `json:"nextblockhash"`
}

type InvokeFunctionStackResult struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

type TokenBalanceResult struct {
	Script      string                      `json:"script"`
	State       string                      `json:"state"`
	GasConsumed string                      `json:"gas_consumed"`
	Stack       []InvokeFunctionStackResult `json:"stack"`
	Tx          string                      `json:"tx"`
}
