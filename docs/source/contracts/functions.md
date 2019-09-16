# Functions

Smart contracts access and modify state variables through functions. Function can be modified with PAYABLE, CONSTANT, EXPORT. PAYABLE, CONSTANT, EXPORT need to be written before the return value of function. The modified function is externally accessible, and the unmodified function is an internal function that is not accessed externally.

## Init Function

The init function is required and must be a function that can be accessed externally. A smart contract only allows one init function. The init function allows arguments to be passed  and will only executed once during the contract deployment and run.

```c++
#include "dipc/dipc.h"

class YourContractName : public Contract {
     EXPORT void init();
}

```

## Accessible Functions
The parameter types and return value types of functions that can be accessed externally are restricted. Currently only simple types are supported:

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

Input parameters and return values do not support storage types and other custom types. Accessible Functions of the same name are not supported, ie overloads of accessible functions are not supported.
There are three types of accessible functions: CONSTANT, PAYABLE, and EXPORT. The usage and functions of the CONSTANT, PAYABLE, and EXPORT macro definitions in dipc are:

| Macro | Utilized Location | Functions | 示例 |
| --- | --- | --- | --- |
|EXPORT | Before the return value of the method declaration and definition | Indicates that the method is an external method | EXPORT  void init(); |
|CONSTANT | The same as above | Indicates that the method does not change the state of the contract data, and can be called directly without sending a transaction. | CONSTANT uint64 getBalance(string addr); |
|PAYABLE | The same as above | Indicates that DIP can be transferred to the contract account by this method. | PAYABLE void transfer(string toAddr, uint_64 value); |
These three macros are independent of each other and cannot be used at the same time.

## Internal Functions
The internal function follows the definition of the C++ language function and does not impose any restrictions.

## Return Value Display

If there is a query request for the return value of the externally accessible function, you need to manually call the DIPC_EMIT_EVENT in the EVENT section in the contract to save it into the Log in the receipts, and then query it by getting the receipts or Log.

## Standard Library Functions

| Function Name  | Function Introduction                                        | Parameters                     | Return Types |
| -------------- | ------------------------------------------------------------ | ------------------------------ | ------------ |
| gasPrice       | Get the gas price of the current transaction                 |                                | int64_t      |
| blockHash      | Get the hash of the block based on the block height          | int64_t number                 | h256         |
| number         | Get the blocknumber of the current block                     |                                | uint64_t     |
| gasLimit       | Get the gas limit of the current transaction                 |                                | uint64_t     |
| timestamp      | Get the packing timestamp of the block                       |                                | address      |
| coinbase       | Get the packaged miner address of the current block          |                                | string       |
| balance        | Get the account balance of an account on the chain           | Address adr                    | uint64       |
| origin         | Get the account address of the contract creator              |                                | Address      |
| caller         | Get the account address of the contract caller               |                                | Address      |
| sha3           | Sha3 encryption operation                                    |                                | h256         |
| getCallerNonce | Get the transaction nonce of the contract caller account     |                                | string       |
| callTransfer   | Transfer the DIP of the contract account to the specified account | Address to ,u256 value         | int64_t      |
| prints         | Print a string variable                                      | string                         | void         |
| prints_l       | Print the first few characters of a string variable          | bool condition, string msg     | void         |
| printi         | Print a 64-bit signed Integer                                | string msg                     | void         |
| printui        | Print a 64-bit unsigned Integer                              | bool condition, string msg     | void         |
| printi128      | Print a 128-bit signed Integer                               | (address addr, uint256 amount) | void         |
| printui128     | Print a 128-bit unsigned Integer                             | const uint128_t* value         | bool         |
| printhex       | Print data in hexadecimal format                             | int64 value                    | string       |
| print          | Template function to print any basic type data               | any basic type                 | void         |
| println        | Template function, print any basic type data, and add a newline at the end | any basic type                 | void         |
| DipcAssert     | Determine if the given condition is true, if it is not true, it will throw an exception | uint64 value                   |              |
| DipcAssertEQ   | Determine if two conditions are equal, and throw an exception if they are not equal | string value                   |              |
| DipcAssertNE   | Determine if two conditions are not equal, and throw an exception if they are equal | two arbitrary expressions      |              |