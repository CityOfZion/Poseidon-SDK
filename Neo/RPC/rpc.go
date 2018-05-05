package rpc

import (
	"github.com/ybbus/jsonrpc"
)

// This package does two things and only two things:
// - connect to the best main, test or private node
// - Use RPC commands t send and recieve data.
// No composition of transactions, no sorting data.
// If it is an argument for the rpc command, it is ready to be sent off
// right now, we only want the RPC to use two commands:
// - Send Transaction and GetRawTransaction
// This package will be called by other packages such as transaction

type rpc struct {
	url string
}

func (u *rpc) SendTransaction(raw string) bool {
	rpcClient := jsonrpc.NewClient(u.url)
	response, err := rpcClient.Call("sendrawtransaction", raw)
	if err != nil {
		return false
	}

	result, _ := response.GetBool()

	return result
}

func (u *rpc) GetRawTransaction(txid string, verbose int8) (string, error) {
	rpcClient := jsonrpc.NewClient(u.url)
	response, err := rpcClient.Call("getrawtransaction", txid, verbose)
	if err != nil {
		return "", err
	}

	return response.GetString()
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
