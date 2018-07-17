package api

// REFER TO README

// import (
// 	"encoding/json"
// 	"fmt"
// 	network "multicrypt/network"
// )

// type API struct{}

// type Unspent struct {
// 	Value float64 `json:"value"`
// 	Txid  string  `json:"txid"`
// 	Index uint32  `json:"n"`
// }

// type JsonBalanceResponse struct {
// 	Balance []struct {
// 		UTXOs  []Unspent `json:"unspent"`
// 		Asset  string    `json:"asset"`
// 		Amount float64   `json:"amount"`
// 	} `json:"balance"`
// 	Address string `json:"address"`
// }

// func (n *API) CheckBalance(address string) JsonBalanceResponse {

// 	url := fmt.Sprintf("https://api.neoscan.io/api/main_net/v1/get_balance/%s", address)

// 	body, err := network.MakeRequest(url)

// 	if err != nil {
// 		return JsonBalanceResponse{}
// 	}

// 	jsonBal := JsonBalanceResponse{}

// 	jsonErr := json.Unmarshal(body, &jsonBal)

// 	if jsonErr != nil {
// 		return JsonBalanceResponse{}
// 	}
// 	// fmt.Println(jsonBal.Balance[3])

// 	for _, bal := range jsonBal.Balance {

// 		fmt.Println(bal.Asset, bal.Amount)
// 	}
// 	return jsonBal
// }

// func (n *API) SendTransaction(b []byte) {

// }

// func (*API) ListAllTokens() []byte {
// 	return []byte("")
// }
