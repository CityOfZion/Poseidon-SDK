package hello

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	neo "multicrypt/Neo"
	network "multicrypt/Neo/RPC"
	cl "multicrypt/Neo/Transactions/ClaimTransaction"
	ct "multicrypt/Neo/Transactions/ContractTransaction"
	ci "multicrypt/Neo/Transactions/InvocationTransaction"
	"strings"
	"time"

	fixed8 "multicrypt/Utils/Fixed8"
	"strconv"
)

// This class could be cleaned up, there is not really a standard for gomobile sdk's
var Coin *neo.Coin

var Contract ct.ContractTransaction

// This class will act as a wrapper, so that the binder does not complain about types

func NewAccount() string {
	return Coin.GenerateMnemonic()
}

func GenerateAddress(mnemonic string, external, internal int) string {
	prKey := Coin.GeneratePrivateKey(mnemonic, external, internal)
	address := Coin.PubKeyToAddress(prKey.PublicKey())
	return address
}

// Prepare refers to converting it from a json to a golang struct
func PrepareTransactionAndSign(transType int, transInJson []byte, keysJsonData []byte, mnemonic string) string {
	fmt.Println(transType)
	if transType == 128 {
		return prepareContractTransactionAndSign(transInJson, mnemonic)
	} else if transType == 2 {
		return prepareClaimTransactionAndSign(transInJson, keysJsonData, mnemonic)
	} else if transType == 209 {

	}
	return prepareContractTransactionAndSign(transInJson, mnemonic)

}

// Due to the script having any type and Swift not allowing me to encode a type of any. We will have seperate function for invocation
// Ideally would want a better solution
func PrepareInvocationTransactionAndSign(operation string, transInJson []byte, scriptInJson []byte, keysJsonData []byte, mnemonic string) string {
	tx := ci.CreateInvocationTransaction()
	json.Unmarshal(transInJson, &tx)
	if operation == "transfer" {
		var s scriptB
		json.Unmarshal(scriptInJson, &s)
		from := s.Args[0].(string)
		to := s.Args[1].(string)
		tx.AppendScriptAttribute(from)
		amount := int64(s.Args[2].(float64)) // when int is passed through, it is converted to float64
		// TODO: test this for large numbers, Alt: send through only strings and use strconv
		fmt.Println("amount in go", amount)
		tx.Script.SetScriptHash(s.Scripthash)
		tx.TransferTokens(from, to, amount)

	}
	var keys Keys
	json.Unmarshal(keysJsonData, &keys)

	for _, k := range keys.Keys {
		key := strings.Split(k, "-")
		External, _ := strconv.Atoi(key[0])
		Index, _ := strconv.Atoi(key[1])
		tx.AddPrivateKey(mnemonic, External, Index)
	}
	tx.Script.SetOperation(operation)

	tx.AddAttributes(uint16(0xff), []byte(time.Now().String()))
	tx.SignTransaction(mnemonic)
	buf := new(bytes.Buffer)
	err := tx.SerialiseTransaction(buf, true)
	if err != nil {
		return ""
	}
	return tx.GetHash() + hex.EncodeToString(buf.Bytes())
}

type scriptB struct {
	Scripthash string
	Operation  string
	Args       []interface{}
}

type Keys struct {
	Keys []string
}

// TODO: refactor, let transactions conform to an interface then
// unmarshall in if statement and pass them to one function to sign
func prepareContractTransactionAndSign(transInJson []byte, mnemonic string) string {
	tx := ct.CreateContractTransaction()
	json.Unmarshal(transInJson, &tx)

	tx.SignTransaction(mnemonic)
	buf := new(bytes.Buffer)
	err := tx.SerialiseTransaction(buf, true)
	if err != nil {

		return err.Error()
	}
	return tx.GetHash() + hex.EncodeToString(buf.Bytes())
}
func prepareClaimTransactionAndSign(transInJson []byte, keysJsonData []byte, mnemonic string) string {
	tx := cl.CreateClaimTransaction()
	json.Unmarshal(transInJson, &tx)

	var keys Keys
	json.Unmarshal(keysJsonData, &keys)

	for _, k := range keys.Keys {
		key := strings.Split(k, "-")
		External, _ := strconv.Atoi(key[0])
		Index, _ := strconv.Atoi(key[1])
		tx.AddPrivateKey(mnemonic, External, Index)

	}

	tx.SignTransaction(mnemonic)
	buf := new(bytes.Buffer)
	err := tx.SerialiseTransaction(buf, true)
	if err != nil {
		return ""
	}
	return tx.GetHash() + hex.EncodeToString(buf.Bytes())
}
func StringToFixed8(s string) string {
	num, err := fixed8.Fixed8DecodeString(s)
	if err != nil {
		return "0"
	}

	return strconv.FormatInt(int64(num), 10)

}
func SendTransaction(url string, transactionData string) bool {
	client := network.Rpc{}
	client.Url = url
	res := client.SendTransaction(transactionData)
	return res
}
