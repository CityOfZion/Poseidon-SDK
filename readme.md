# ReadMe

This repo provides all of the necessary functions to run a BIP32, BIP44 and BIP39 compliant wallet on NEO.

The framework also allows for more than one cryptocurrency, with different elliptic curves. Please see the BTC package on how to use this package with another cryptocurency.

# Usability

This package is usable as of V 0.30.0 however this repo should be forked as there will be breaking changes until V 1.0


# SDK API

## Key Management

### Generate a mnemonic

    neoCoin := Coin{}
    seed := neoCoin.GenerateMnemonic()

### Generate Private Key

    neoCoin := Coin{}
	mnemonic := "online ramp onion faculty trap clerk near rabbit busy gravity prize employ"
	privKey := neoCoin.GeneratePrivateKey(mnemonic, 0, 0) // returns a bip32.Key
    hex.EncodeToString(privKey.Key)

### Generate WIF From Private Key

    neoCoin := Coin{}
	mnemonic := "online ramp onion faculty trap clerk near rabbit busy gravity prize employ"
	privKey := neoCoin.GeneratePrivateKey(mnemonic, 0, 1)
    crypto.WIFEncode(privKey.Key)

### Generate NEO Address From Private Key

    neoCoin := Coin{}
	mnemonic := "online ramp onion faculty trap clerk near rabbit busy gravity prize employ"
	prKey := neoCoin.GeneratePrivateKey(mnemonic, 0, 1)
	address := neoCoin.PubKeyToAddress(prKey.PublicKey())

### Generate NEO Address From Public Key

    neoCoin := Coin{}
	mnemonic := "online ramp onion faculty trap clerk near rabbit busy gravity prize employ"
	privKey := neoCoin.GeneratePrivateKey(mnemonic, 0, 1) // returns a bip32.Key
	pubKey := privKey.PublicKey()
	address := neoCoin.PubKeyToAddress(pubKey)

## Script Builder

### Generate Script For Given Arguments

    s := Script{
		Operation:  "name",
		Args:       []interface{}{},
		ScriptHash: "5b7074e873973a6ed3708862f219a6fbf4d1c411",
	}
    s.Hex()

*Please note, that the argument slice will not accept structs.*

## Transactions

### Contract Transaction

*To send 1000 GAS to ANmQW9femWe1oPgEDD4pSTaNUu5rvHYW1R and 1000 GAS to ALxUuc4gSNaYdUUrPudem4DMkowz6x6Rwo*

    tx := CreateContractTransaction()

	tx.AddOutput(GAS, int64(1000 * 1e8), "ANmQW9femWe1oPgEDD4pSTaNUu5rvHYW1R")

	tx.AddOutput(GAS, int64(1000 * 1e8), "ALxUuc4gSNaYdUUrPudem4DMkowz6x6Rwo")

	tx.AddInput("f08dfee4598e3a91e20b151ad63a967b95240c5ec96ecf4548e27dcff7a330d7", 0, "ALxUuc4gSNaYdUUrPudem4DMkowz6x6Rwo", 0, 1)

	tx.SignTransaction("online ramp onion faculty trap clerk near rabbit busy gravity prize employ")

	buf := new(bytes.Buffer)

	err := tx.SerialiseTransaction(buf, true)

    transactionInHex := hex.EncodeToString(buf.Bytes())

*Please note that the amount needs to be in Fixed8 format. Numbers are expected to be in their Fixed8 format as an int64. E.g. 5 GAS would be int64(5 * 1e8) . Users can use the Fixed8 class to automatically do the conversion, or have this done in their applications business logic.*

### Invocation Transaction

	tx := CreateInvocationTransaction()

	tx.AppendScriptAttribute("ALxUuc4gSNaYdUUrPudem4DMkowz6x6Rwo")
	tx.Script.SetScriptHash("9aff1e08aea2048a26a3d2ddbb3df495b932b1e7")

	tx.TransferTokens("ALxUuc4gSNaYdUUrPudem4DMkowz6x6Rwo", "AQGV3Q4Q4kczcb7Q6Ax11zBP5WcygRRKXw", 100)
	tx.AddAttributes(uint16(0xff), []byte(time.Now().String()))
	tx.AddPrivateKey("online ramp onion faculty trap clerk near rabbit busy gravity prize employ", 0, 1)
	tx.SignTransaction("online ramp onion faculty trap clerk near rabbit busy gravity prize employ")
	buf := new(bytes.Buffer)
	err := tx.SerialiseTransaction(buf, true)
    transactionInHex := hex.EncodeToString(buf.Bytes())

### Claim Transaction

    tx := CreateClaimTransaction()
	tx.AddClaim("480a72a7d43452f9503f4f9552768cfd54ffa103776b703e7519d22d959ade68", 0, "ALxUuc4gSNaYdUUrPudem4DMkowz6x6Rwo", 0, 1)
	val, _ := fixed8.Fixed8DecodeString("0.02119752")
	tx.AddOutput(GAS, int64(val), "ALxUuc4gSNaYdUUrPudem4DMkowz6x6Rwo")
	tx.AddPrivateKey("online ramp onion faculty trap clerk near rabbit busy gravity prize employ", 0, 1)
	tx.SignTransaction("online ramp onion faculty trap clerk near rabbit busy gravity prize employ")
	buf := new(bytes.Buffer)
	err := tx.SerialiseTransaction(buf, true)
    transactionInHex := hex.EncodeToString(buf.Bytes())

## RPC

### Send Transaction

	node := Rpc{"http://seed1.cityofzion.io:8080"}
	res := node.SendTransaction(transactionInHex)
    *// Transaction successful if res == true*

### Get Transaction

    txid := "56d477f7cffe5f3f919be798e7c752c754faf03870d202432b8157d9da0dfc57"

    node := Rpc{"http://seed3.aphelion-neo.com:10332"}
	result, err := node.GetRawTransaction(txid, 0)
    *// 0 for transaction in hex, 1 for JSON*

### Get Block

    blockNumber := 2000
    node := Rpc{"https://seed1.neo.org:20331"}
	res, err := node.getRawBlock(blockNumber, 0)

### RPC Invoke script 

    script :="1457c4cf51f12ce6d78a585f1ea9bd1f3927c7232c51c10962616c616e63654f6667e7b132b995f43dbbddd2a3268a04a2ae081eff9a"

    node := Rpc{"https://seed1.neo.org:20331"}

	res, err := node.InvokeScript()

*Script generated from script builder package*


