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