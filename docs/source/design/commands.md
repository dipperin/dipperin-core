# Command Line Tool
- [How to use command line tool?](#How-to-use-command-line-tool)
- [How to operate Test Net Node?](#How-to-operate-Test-Node)
- [Related Functional Operations](#Related-Functional-Operations)

## How to use command line tool?

### Connect to test environment

Dipperin command line tools are located in the $GOBIN directory: `~/go/bin/Dipperincli`.

Monitor for test environment: `http://${TestServer}:8887`.

PBFT whole process demonstration: `http://${TestServer}:8888`.

If the startup command line needs to manipulate other startup nodes,
it can be done by specifying parameters in the startup command.

Example:

Assuming that the local cluster is currently started, the command line tool needs to manipulate the V0 node.

 IP: `127.0.0.1`.

HttpPort: `50007`.

```
dipperincli -- http_host 127.0.0.1 --http_port 50007
```

Connect to the test environment:
```
boots_env=test ~/go/bin/dipperincli
```

Or set temporary environment variables first:
```
export boots_env=test
```

### Start Dipperin node

The following command is to start a node, which requires a wallet password.

If no wallet path is specified, the default system path is used: `~/.dipperin/`.

```
dipperincli -- node_type [type] -- soft_wallet_pwd [password]
```

Example:

Local startup miner:
```
dipperincli -- node_type 1 -- soft_wallet_pwd 123
```

Local startup miner(start mining):
```
dipperincli -- node_type 1 -- soft_wallet_pwd 123 -- is_start_mine 1
```

Local startup verifier:
```
dipperincli -- node_type 2 -- soft_wallet_pwd 123
```

Connect to the test environment:
```
boots_env=test ~/go/bin/dipperincli -- soft_wallet_pwd 123
```

### Error

If dipperincli started in a wrong way,
it may be that the local link data is not synchronized with the link state,
and the local link data needs to be deleted:
```
cd ~
rm .dipperin -fr
```

restart command line tool


## Related Functional Operations

Separate multiple parameters by','
```
moduleName MethodName 【-p parameters]

```

### Block chain
{Text: "GetReceiptsByBlockNum", Description: ""},

Get current block:
```
chain CurrentBlock
```

Get genesis block:
```
chain GetGenesis
```

Get block by number
```
chain GetBlockByNumber -p blockNumber
chain GetBlockByNumber -p 1
```

Get block by block hash:
```
chain GetBlockByHash -p blockHash
chain GetBlockByHash -p 0x0f7057ff3e3048ed38c0ac2353e001dad6aded5d825d43fcc924a39221713e4c
```

Get Peers
```
chain GetPeers
```

Add Peer
```
chain AddPeer -p url
chain AddPeer -p enode://199cc6526cb63866dfa5dc81aed9952f2002b677560b6f3dc2a6a34a5576216f0ca25711c5b4268444fdef5fee4a01a669af90fd5b6049b2a5272b39c466b2ac@127.0.0.1:10006
```

### Miner master

Start mining:
```
miner StartMine
```

Stop mining:
```
miner StopMine
```

Set miner address:
```
miner SetMineCoinBase -p address
miner SetMineCoinBase -p 0x0000e447B8B7851D3FBD5C6A03625D288cfE9Bb5eF0E
```

### Personal

#### Wallet

Look up local wallet:
```
personal ListWallet
```

Look up local wallet account:

If the wallet type and path are not specified, the default wallet is displayed
```
personal ListWalletAccount [-p walletType,walletPath]
personal ListWalletAccount -p SoftWallet,~/tmp/dipperin_apps/default_v0/CSWallet
```

Create new wallet:
```
personal EstablishWallet -p walletType,walletPath,password
personal EstablishWallet -p SoftWallet,/tmp/TestWallet,123

```

Recovery wallet:
```
personal RestoreWallet -p walletType,walletPath,password,passpharse,mnemonic,...,mnemonic
personal RestoreWallet -p SoftWallet,/tmp/TestWallet2,123,,plastic,balcony,trophy,fuel,vacant,inmate,profit,rival,mimic,cute,hurdle,pig,column,pudding,visit,edge,rhythm,armed,cook,federal,amount,stock,damp,bring
```

Open wallet:

If the wallet type and path are not specified, the default wallet is displayed
```
personal OpenWallet -p walletType,walletPath,password
personal OpenWallet -p SoftWallet,/tmp/TestWallet3,123
```

Close wallet:

If the wallet type and path are not specified, the default wallet is displayed
```
personal CloseWallet [-p walletType,walletPath]
personal CloseWallet -p SoftWallet,/tmp/TestWallet3
```
	
#### Account

Add account:

If the wallet type and path are not specified, the default wallet is displayed
```
personal AddAccount -p walletType,walletPath
personal AddAccount -p SoftWallet,/tmp/TestWallet3
```

Get account current balance:
```
personal CurrentBalance [-p address]
personal CurrentBalance -p 0x0000e447B8B7851D3FBD5C6A03625D288cfE9Bb5eF0E
```

Get account deposit:
```
personal CurrentStake [-p address]
personal CurrentStake -p 0x0000e447B8B7851D3FBD5C6A03625D288cfE9Bb5eF0E
```

Get account nonce:
```
personal GetAddressNonceFromWallet [-p address]
personal GetAddressNonceFromWallet -p 0x00001c2beC8E0E4caac668cD75d520E41f827092Ce79
```

Get account Reputation:
```
personal CurrentReputation [-p address]
personal CurrentReputation -p 0x00001c2beC8E0E4caac668cD75d520E41f827092Ce79
```

Get Default Account Balance
```
personal GetDefaultAccountBalance
```

Get Default Account Stake
```
personal GetDefaultAccountStake
```

Get Transaction Nonce
```
personal GetTransactionNonce [-p address]
personal GetTransactionNonce -p 0x000062be10f46b5d01Ecd9b502c4bA3d6131f6fc2e41
```

Set wallet signer（default account）
```
personal SetBftSigner -p address
personal SetBftSigner -p 0x000062be10f46b5d01Ecd9b502c4bA3d6131f6fc2e41
```

### Transaction

Register verifier:
```
tx SendRegisterTransaction -p from,deposit,transactionFee
tx SendRegisterTransaction -p 0x00004179D57e45Cb3b54D6FAEF69e746bf240E287978,100,11

or

tx SendRegisterTx -p deposit,transactionFee
tx SendRegisterTx -p 100,11
```

Unregister verifier:
```
tx SendCancelTransaction -p from,transactionFee
tx SendCancelTransaction -p 0x00004179D57e45Cb3b54D6FAEF69e746bf240E287978,11

or 

tx SendCancelTx -p transactionFee
tx SendCancelTx -p 11
```

Redemption of the deposit:
```
tx SendUnStakeTransaction -p from,transactionFee
tx SendUnStakeTransaction -p 0x00004179D57e45Cb3b54D6FAEF69e746bf240E287978,11

or 

tx SendUnStakeTx -p transactionFee
tx SendUnStakeTx -p 11
```

Get transation nonce:
```
tx GetTransactionNonce -p [address]
tx GetTransactionNonce -p 0x00001c2beC8E0E4caac668cD75d520E41f827092Ce79
```

Send normal transaction:
```
tx SendTransaction -p from,to,value,transactionFee,[extradata]
tx SendTransaction -p 0x00004179D57e45Cb3b54D6FAEF69e746bf240E287978,0x00007eDe4D5D808DA8a267284b38E00ABccb42889dF2,20000,10

or 

tx SendTx -p to,value,transactionFee,[extradata]
tx SendTx -p 0x00007eDe4D5D808DA8a267284b38E00ABccb42889dF2,20000,10
```

Get transaction
```
tx Transaction [TxId]
tx Transaction -p 0xf8dd21db65b2adcb5e3ed3c61475eb66a1653d309b1a82354959fdf58852f023
```

Send create contract transaction
```
tx SendTransactionContract -p from,to,gasLimit,[gasPrice] --abi ~/dipc/cmake-build-debug/token/token.cpp.abi.json --wasm ~/dipc/cmake-build-debug/token/token5.wasm [--input pramas...] --is-create
tx SendTransactionContract -p 0x000062be10f46b5d01Ecd9b502c4bA3d6131f6fc2e41,10,1,1 -abi ~/dipc/cmake-build-debug/token/token.cpp.abi.json -wasm ~/dipc/cmake-build-debug/token/token5.wasm -input dipp,DIpp,10000000 -isCreate
```

Send call contract transaction
```
tx SendTransactionContract -p from,to,gasLimit,[gasPrice] --func-name funcName [--input pramas...]
tx SendTransactionContract -p 0x000062be10f46b5d01Ecd9b502c4bA3d6131f6fc2e41,10,1,1 --func-name getBalance --input 0x000062be10f46b5d01Ecd9b502c4bA3d6131f6fc2e41
```

Get Contract Address By TxHash
```
tx GetContractAddressByTxHash -p txHash
tx GetContractAddressByTxHash -p 0x805bb426f32382cf8731f03aaaec7599bf3d57b1548efd8f6a87ce49b54357d4
```

Get Convert Receipt By TxHash
```
tx GetConvertReceiptByTxHash -p txHash
tx GetConvertReceiptByTxHash -p 0x805bb426f32382cf8731f03aaaec7599bf3d57b1548efd8f6a87ce49b54357d4
```

Get Receipt By TxHash
```
tx GetReceiptByTxHash -p txHash
tx GetReceiptByTxHash -p tx GetConvertReceiptByTxHash -p 0x805bb426f32382cf8731f03aaaec7599bf3d57b1548efd8f6a87ce49b54357d4
```

Estimate Gas for contract transaction create
```
tx EstimateGas -p from,to,gasLimit,[gasPrice] --abi ~/dipc/cmake-build-debug/token/token.cpp.abi.json --wasm ~/dipc/cmake-build-debug/token/token5.wasm [--input pramas...] --is-create
tx EstimateGas -p 0x000062be10f46b5d01Ecd9b502c4bA3d6131f6fc2e41,10,1,1 -abi ~/dipc/cmake-build-debug/token/token.cpp.abi.json -wasm ~/dipc/cmake-build-debug/token/token5.wasm -input dipp,DIpp,10000000 -isCreate
```

Estimate Gas for contract transaction call
```
tx EstimateGas -p from,to,gasLimit,[gasPrice] --func-name funcName [--input pramas...]
tx EstimateGas -p 0x000062be10f46b5d01Ecd9b502c4bA3d6131f6fc2e41,10,1,1 --func-name getBalance --input 0x000062be10f46b5d01Ecd9b502c4bA3d6131f6fc2e41
```

#### ERC20

Create ERC20 contract:

The contract address must be saved for successful creation
```
tx AnnounceERC20 -p owner_address, token_name, token_symbol, token_total_supply,decimal,transactionFee
tx AnnounceERC20 -p 0x0000D07252C7A396Cc444DC0196A8b43c1A4B6c53532,chain,stack,5,3,0.00001
```

Look up ERC20 information:
```
tx ERC20GetInfo -p contract_address
tx ERC20GetInfo -p 0x00100f35adf022a8aaAbef59abB97665788CDdbA30e3
```

Look up ERC20 total supply
```
tx ERC20TotalSupply -p contract_address
tx ERC20TotalSupply -p 0x00100f35adf022a8aaAbef59abB97665788CDdbA30e3
```

Look up ERC20 token name (No need) :
```
tx ERC20TokenName -p contract_address
tx ERC20TokenName -p 0x00100f35adf022a8aaAbef59abB97665788CDdbA30e3
```

Look up ERC20 token symbol (No need) :
```
tx ERC20TokenSymbol -p contract_address
tx ERC20TokenSymbol -p 0x00100f35adf022a8aaAbef59abB97665788CDdbA30e3
```

Look up ERC20 token decimals (No need) :
```
tx ERC20TokenDecimals -p contract_address
tx ERC20TokenDecimals -p 0x00100f35adf022a8aaAbef59abB97665788CDdbA30e3
```

Look up ERC20 token quota:
```
tx ERC20Allowance -p contract_address,from_address,to_address
tx ERC20Allowance -p 0x00100f35adf022a8aaAbef59abB97665788CDdbA30e3,0x0000D07252C7A396Cc444DC0196A8b43c1A4B6c53532,0x0000B04985A7ccc00ab023d9bC40E241F9DF0379d8c4
```

Allocation of ERC20 contract quota:
```
tx ERC20Approve -p contract_address,from_address,to_address,amount,transactionFee
tx ERC20Approve -p 0x00100f35adf022a8aaAbef59abB97665788CDdbA30e3,0x0000D07252C7A396Cc444DC0196A8b43c1A4B6c53532,
0x0000B04985A7ccc00ab023d9bC40E241F9DF0379d8c4,4,0.00001
```

Look up ERC20 balance:
```
tx ERC20Balance -p contract_address,owner_address
tx ERC20Balance -p 0x00100f35adf022a8aaAbef59abB97665788CDdbA30e3,0x0000D07252C7A396Cc444DC0196A8b43c1A4B6c53532
```

ERC20 contract transfer:
```
tx ERC20Transfer -p contract_address,from_address,to_address,amount,transactionFee
tx ERC20Transfer -p 0x00100f35adf022a8aaAbef59abB97665788CDdbA30e3,0x0000D07252C7A396Cc444DC0196A8b43c1A4B6c53532,
0x0000B04985A7ccc00ab023d9bC40E241F9DF0379d8c4,4,0.00001
```

ERC20 contract transfer from:
```
tx ERC20TransferFrom -p contract_address,owner_address,from_address,to_address,amount,transactionFee
tx ERC20TransferFrom -p 0x0000B04985A7ccc00ab023d9bC40E241F9DF0379d8c4,0x0000D07252C7A396Cc444DC0196A8b43c1A4B6c5353d,
0x0000D07252C7A396Cc444DC0196A8b43c1A4B6c53532,0x00100f35adf022a8aaAbef59abB97665788CDdbA30e3,4,0.00001
```

Transfer EDIP To DIP
```
tx TransferEDIPToDIP -p from,eDIPValue,transactionFee
tx TransferEDIPToDIP -p 0x0000B04985A7ccc00ab023d9bC40E241F9DF0379d8c4,10000,100
```

Set Exchange Rate
```
tx SetExchangeRate -p from,exchangeRate,transactionFee
tx SetExchangeRate -p 0x0000B04985A7ccc00ab023d9bC40E241F9DF0379d8c4,0.5,100
```



### Verifiers
	{Text: "CheckVerifierType", Description: ""},

Get Verifiers by slot:
```
verifier GetVerifiersBySlot -p [round]
verifier GetVerifiersBySlot -p 2

```

Get current slot verifiers:
```
verifier GetCurVerifiers
```

Get next slot verifiers:
```
verifier GetNextVerifiers
```

Get verifier status
```
verifier VerifierStatus [-p accountAddress]
```

Get Block Different Verifier Info
```
verifier GetBlockDiffVerifierInfo -p blockNumber
verifier GetBlockDiffVerifierInfo -p 11
```

Check An Address Is Or Not Verifier In Some Slot
```
verifier CheckVerifierType -p slot,address
verifier CheckVerifierType -p 11,0x000062be10f46b5d01Ecd9b502c4bA3d6131f6fc2e41
```



