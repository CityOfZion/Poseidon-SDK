package rpc

import (
	"fmt"

	"github.com/ybbus/jsonrpc"
)

type Rpc struct {
	Url string
}

func (u *Rpc) SendTransaction(raw string) bool {
	rpcClient := jsonrpc.NewClient(u.Url)
	response, err := rpcClient.Call("sendrawtransaction", raw)
	if err != nil {
		return false
	}

	result, _ := response.GetBool()
	return result
}

func (u *Rpc) GetRawTransaction(txid string, verbose int8) (string, error) {
	rpcClient := jsonrpc.NewClient(u.Url)
	response, err := rpcClient.Call("getrawtransaction", txid, verbose)
	if err != nil {
		return "", err
	}

	return response.GetString()
}

func (u *Rpc) InvokeScript(script string) (TokenBalanceResult, error) {
	rpcClient := jsonrpc.NewClient(u.Url)
	response, err := rpcClient.Call("invokescript", script)
	if err != nil {
		return TokenBalanceResult{}, err
	}
	var res *TokenBalanceResult
	err = response.GetObject(&res)
	if err != nil || res == nil {
		return TokenBalanceResult{}, err
	}
	return *res, nil
}

func (u *Rpc) GetRawMempool() []string {
	rpcClient := jsonrpc.NewClient(u.Url)

	response, err := rpcClient.Call("getrawmempool", "")
	res := []string{}
	if err != nil || response == nil {
		return res
	}

	transactions := response.Result.([]interface{})

	for _, transaction := range transactions {
		if transactionAsString, ok := transaction.(string); ok {
			res = append(res, transactionAsString)
		}

	}

	return res
}

func (u *Rpc) GetBlock(index int) (BlockRes, error) {
	rpcClient := jsonrpc.NewClient(u.Url)
	response, err := rpcClient.Call("getblock", index, 1)
	fmt.Println(response)
	if err != nil {

		return BlockRes{}, err
	}

	var res *BlockRes

	err = response.GetObject(&res) // expects a rpc-object result value like: {"id": 123, "name": "alex", "age": 33}
	if err != nil || res == nil {
		// some error on json unmarshal level or json result field was null

		return BlockRes{}, err
	}
	return *res, nil

}
func (u *Rpc) getRawBlock(index int, verbose int) (string, error) {
	rpcClient := jsonrpc.NewClient(u.Url)
	response, err := rpcClient.Call("getblock", index, verbose)
	fmt.Println(response)
	if err != nil {

		return "", err
	}

	res, err := response.GetString() // expects a rpc-object result value like: {"id": 123, "name": "alex", "age": 33}
	return res, err

}

// TODO: I would like to start off with a list of hardcoded nodes,
//and from there we use getpeers command to build up a list of good nodes
// However: When I get a peer and try to connect to them.
// for some reason, they do not respond
// For now, we will just get a node from:

func getAvailableNode() string {
	return "http://seed1.cityofzion.io:8080"
}

/*

“http://seed1.cityofzion.io:8080”,
 “http://seed2.cityofzion.io:8080”,
  “http://seed3.cityofzion.io:8080”, “http://seed4.cityofzion.io:8080”,
   “http://seed5.cityofzion.io:8080”, “http://api.otcgo.cn:10332”,
    “https://seed1.neo.org:10331”, “http://seed2.neo.org:10332”,
   “http://seed3.neo.org:10332”, “http://seed4.neo.org:10332”,
  “http://seed5.neo.org:10332”, “http://seed0.bridgeprotocol.io:10332”,
  “http://seed1.bridgeprotocol.io:10332”, “http://seed2.bridgeprotocol.io:10332”, “http://seed3.bridgeprotocol.io:10332”,
   “http://seed4.bridgeprotocol.io:10332”, “http://seed5.bridgeprotocol.io:10332”,
  “http://seed6.bridgeprotocol.io:10332”, “http://seed7.bridgeprotocol.io:10332”,
  “http://seed8.bridgeprotocol.io:10332”, “http://seed9.bridgeprotocol.io:10332”,
  “http://seed1.redpulse.com:10332”, “http://seed2.redpulse.com:10332”,
  “https://seed1.redpulse.com:10331”, “https://seed2.redpulse.com:10331”, “http://seed1.treatail.com:10332”,
  “http://seed2.treatail.com:10332”, “http://seed3.treatail.com:10332”,
  “http://seed4.treatail.com:10332”, “http://seed1.o3node.org:10332”, “http://seed2.o3node.org:10332”,
  “http://54.66.154.140:10332”, “http://seed1.eu-central-1.fiatpeg.com:10332”,
  “http://seed1.eu-west-2.fiatpeg.com:10332”, “http://seed1.aphelion.org:10332”,
  “http://seed2.aphelion.org:10332”, “http://seed3.aphelion.org:10332”,
  “http://seed4.aphelion.org:10332”

*/
