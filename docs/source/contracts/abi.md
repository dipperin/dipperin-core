## Understanding ABI Files

When publishing a smart contract on the dipperin chain, you need to provide the ABI file generated when compiling the smart contract using the dipc tool.  ABI (Application Binary Interface) is a JSON-based description that shows how to translate user operations between JSON and binary representations. ABI also describes how to convert database state to JSON or convert from JSON. Once you describe your smart contract through ABI, developers and users can seamlessly interact with your smart contract via JSON.

Special Note: ABI can be bypassed when executing a contract transaction. The messages and actions passed to the smart contract do not have to comply with the ABI. ABI is a guide, not a guard.

All methods that can be called directly by the user in the contract will be described by generating a corresponding JSON object in the ABI file.

```json
[{
    "name": "init",
    "inputs": [
        {
            "name": "tokenName",
            "type": "string"
        },
        {
            "name": "symbol",
            "type": "string"
        },
        {
            "name": "supply",
            "type": "uint64"
        }
    ],
    "outputs": [],
    "constant": "false",
    "payable": "false",
    "type": "function"
},
{
    "name": "GetBalance",
    "inputs": [
        {
            "type": "string"
        },
        {
            "type": "string"
        },
        {
            "type": "uint64"
        }
    ],
    "type": "event"
}
]
```
This is part of an ABI file for an example token contract. The meanings of their fields are:

name:               indicates the name of the method in the contract or the name of the event in the contract;
inputs:              method parameters；
inputs.type:     indicates the type of the input parameter;
inputs.name:   indicates the field name of the input parameter;
outputs:           the return value of the method;
outputs.type:  indicates the type of the return value;
constant:         a value of true means that the method does not change the state of the contract data, and can be called directly without sending a transaction;
payable:          a value of true indicates that DIP can be transferred to the contract account by this method.
type:                indicates the type of the abi object, which has two types: event and function.

The types supported by inputs.type and outputs.type are as follows ( the types of input and return values supported in accessible functions ) : 

  - std::string   
  - unsigned char
  - char[]
  - char *
  - char
  - const char*
  - bool
  - unsigned long long
  - unsigned long
  - unsigned __int128
  - uint128_t
  - uint64_t
  - uint32_t
  - unsigned short
  - uint16_t
  - uint8_t
  - __int128
  - int128_t
  - long long
  - int64_t
  - long
  - int32_t
  - short
  - int16_t
  - int8_t
  - int
