Json Transaction schema for neo
Types:

InvocationTransaction TransactionType = 0xd1

// Not quite sure on how to form the claim transaction
ClaimTransaction TransactionType = 0x02
{
    version:"2.0", // jsonrpc command
    type: "ClaimTransaction",
    SpecificData: "nil",
    Attributes : "nil",
    Claims : [
        {

        }
    ],
    Input : [],
    Output :[],
    Scripts :[]
}



ContractTransaction   TransactionType = 0x80
{
    // This is used to send Neo and Gas to others
    version: "2.0",
    type: "ContractTransaction",
    "SpecificData" : "nil",
    "Attributes" : "nil",
    inputs: [
        {
            prevHash: "Hash of prev transaction",
            prevIndex: "previous transaction's index"
            path: "derivation path for the address, I want to spend from... 0/2"
        },
        {
            prevHash: "Hash of prev transaction2",
            prevIndex: "previous transaction's index"
            path: "derivation path for the address, I want to spend from... 0/2"
        },

    ],
    outputs: [
        {
            AssetId : "AssetId of the asset you want to send",
            Value: "How much of the asset you want to send",
            ScriptHash: "Address of the person you want to send to"
        },
        {
            AssetId : "AssetId of the asset you want to send",
            Value: "How much of the asset you want to send",
            ScriptHash: "Address of the person you want to send to"
        }
        
    ],
    scripts: [
         "signatures used to validate the above inputs and outputs"
    ]


}


This first response is when I put in arguments for invokescript.

&{  2.0 
    map[
        script:00c1046e616d656711c4d1f4fba619f2628870d36e3a9773e874705b 
        state:FAULT, BREAK 
        gas_consumed:0.011 
        stack:[map[type:Array value:[]] map[type:ByteArray value:6e616d65]]
    ] 
    <nil> 
    0
}

<--Invokation Transaction>
This is without arguments for invoke script.
&{ 2.0 
    map[script:00046e616d656711c4d1f4fba619f2628870d36e3a9773e874705b 
        state:FAULT, BREAK 
        gas_consumed:0.01 
        stack:[map[type:ByteArray value:] map[type:ByteArray value:6e616d65]]
    ] 
    <nil> 
    0
}

// With arguments, the stack is stack:[map[type:Array value:[]] map[type:ByteArray value:6e616d65]]
	// Without arguments it is:      stack:[map[type:ByteArray value:] map[type:ByteArray value:6e616d65]]

    So I set the Value to interface{} for now as a quick fix, as I will not be using the stack for now.


Decision:


- Does the client keep the state of the UTXOs and all addresses with funds?

If so, how?

-- who keeps the list of all of the addresses in use?

client could keep it in realmdb and subscribe to notifications, so that users know when they get a payment
   - This can be done with post notification and dynamodb streams.

- client can check for balances, using rpc call? or api endpoint, or let them just use the utxos that they have. Yeah, just let them use their utxos.

- Push notification, will also notify client, so that the app updates in the background.

- Client looks at cached UTXOs before requesting them from api. (This will have to be done in the native language; Swift)
The Go library should therefore only seek out the UTXOs and gas claims for a single address and value or multiple addresses

endpointurl.com/address/<address>?amount=10 This will return all utxos for that address that sum upto or more than amount, which is 10. The utxos will be sorted, with smallest first

endpointurl.com/address/<address> This wll return all utxos for that address

endpointurl.com/claims/<address> This will return all possible claims for address

endpointurl.com/xpub/<xpub> This will be given an xpub and return all utxos for addresses under the xpub

endpointurl.com/xpub/<xpub>?amount=10 This will return a set of utxos that sum to 10 or over from the xpub.


Neon-Js get scriptHash

+const attachAttributesForEmptyTransaction = config => {
+  if (config.tx.inputs.length === 0 && config.tx.outputs.length === 0) {
+    config.tx.addAttribute(TxAttrUsage.Script, reverseHex(getScriptHashFromAddress(config.address)))
+    // This adds some random bits to the transaction to prevent any hash collision.
+    config.tx.addRemark(Date.now().toString() + generatePrivateKey().substr(0, 8))
