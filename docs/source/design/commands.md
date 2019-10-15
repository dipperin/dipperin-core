# Command Line Tool
- [How to use command line tool](#How-to-use-command-line-tool)
- [How to operate Test Node](#How-to-operate-Test-Node)
- [Related Functional Operations](#Related-Functional-Operations)

## How to use command line tool

### Connect to test environment

dipperin command line tools are located in the $GOBIN directory: `~/go/bin/dipperincli`.

Monitor for test environment: `http://10.200.0.139:8887`.

PBFT whole process demonstration: `http://10.200.0.139:8888`.

If the startup command line needs to manipulate other startup nodes,
it can be done by specifying parameters in the startup command.

Example:

Assuming that the local cluster is currently started, the command line tool needs to manipulate the V0 node.

IP: `127.0.0.1`.

HttpPort: `50007`.

```
dipperincli -- http_host 127.0.0.1 -- http_port 50007
```

Connect to the test environment:
```
boots_env = test ~/go/bin/dipperincli
```

Or set temporary environment variables first:
```
export boots_env=test
```

### Start dipperin node

The following command is to start a node, which requires a wallet password.

If no wallet path is specified, the default system path is used: `~/.dipperin/`.

```
dipperincli --node_type [type] --soft_wallet_pwd [password]
```

Example:

Local startup miner:
```
dipperincli --node_type 1 --soft_wallet_pwd 123
```

Local startup miner(start mining):
```
dipperincli --node_type 1 --soft_wallet_pwd 123 --is_start_mine 1
```

Local startup verifier:
```
dipperincli --node_type 2 --soft_wallet_pwd 123
```

Connect to the test environment:
```
boots_env = test ~/go/bin/dipperincli --soft_wallet_pwd 123
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
[ModuleName] [MethodName] -p [parameters]

```

### Transaction methods

AnnounceERC20:
```
tx AnnounceERC20 -p [owner_address],[token_name],[token_symbol],[token_total_supply],[decimal],[gasPrice],[gasLimit]
tx AnnounceERC20 -p 0x0000661A3c6c0955B5E6dbf935f0891aAA1112b9E9ca,wjw,dip,10000,3,10wu,100000
```

ERC20Transfer:
```
tx ERC20Transfer -p [contract_address],[owner],[to_address],[amount],[gasPrice],[gasLimit]
tx ERC20Transfer -p 0x0010Cb4174726E90E3ce09360B5F0488Ab29Fa5aB130,0x0000661A3c6c0955B5E6dbf935f0891aAA1112b9E9ca,0x0000970e8128aB834E8EAC17aB8E3812f010678CF791,1000,10wu,100000
```

ERC20TransferFrom:
```
tx ERC20TransferFrom -p [contract_address],[owner],[from_address],[to_address],[amount],[gasPrice],[gasLimit]
tx ERC20TransferFrom -p 0x0010Cb4174726E90E3ce09360B5F0488Ab29Fa5aB130,0x0000661A3c6c0955B5E6dbf935f0891aAA1112b9E9ca,0x0000661A3c6c0955B5E6dbf935f0891aAA1112b9E9ca,0x0000970e8128aB834E8EAC17aB8E3812f010678CF791,1,10wu,100000
```

ERC20Allowance:
```
tx ERC20Allowance -p [contract_address],[owner],[spender]
tx ERC20Allowance -p 0x0010Cb4174726E90E3ce09360B5F0488Ab29Fa5aB130,0x0000661A3c6c0955B5E6dbf935f0891aAA1112b9E9ca,0x0000970e8128aB834E8EAC17aB8E3812f010678CF791
```

ERC20Approve:
```
tx ERC20Approve -p [contract_address],[owner],[to_address],[amount],[gasPrice],[gasLimit]
tx ERC20Approve -p 0x0010Cb4174726E90E3ce09360B5F0488Ab29Fa5aB130,0x0000661A3c6c0955B5E6dbf935f0891aAA1112b9E9ca,0x0000970e8128aB834E8EAC17aB8E3812f010678CF791,1000,10wu,100000
```

ERC20Balance:
```
tx ERC20Balance -p [contract_address],[owner_address]
tx ERC20Balance -p 0x0010Cb4174726E90E3ce09360B5F0488Ab29Fa5aB130,0x0000970e8128aB834E8EAC17aB8E3812f010678CF791
```

ERC20GetInfo:
```
tx ERC20GetInfo -p [contract_address]
tx ERC20GetInfo -p 0x0010Cb4174726E90E3ce09360B5F0488Ab29Fa5aB130
```

Register verifier:
```
tx SendRegisterTx -p [stake],[gasPrice],[gasLimit]
tx SendRegisterTx -p 1000dip,1wu,21000
```

Unregister verifier:
```
tx SendCancelTx -p [gasPrice],[gasLimit]
tx SendCancelTx -p 1wu,21000
```

Redemption of the deposit:
```
tx SendUnStakeTx -p [gasPrice],[gasLimit]
tx SendUnStakeTx -p 1wu,21000
```

Send transaction:
```
tx SendTx -p [to],[value],[gasPrice],[gasLimit]
tx SendTx -p 0x0000970e8128aB834E8EAC17aB8E3812f010678CF791,100dip,1wu,21000
```

Create contract:
```
tx SendTransactionContract -p [from],[value],[gasPrice],[gasLimit] --abi [abiPath] --wasm [wasmPath] --is-create --input [params]
tx SendTransactionContract -p 0x0000661A3c6c0955B5E6dbf935f0891aAA1112b9E9ca,0dip,1wu,5000000 --abi /home/qydev/testData/token-payable/token-payable.cpp.abi.json --wasm /home/qydev/testData/token-payable/token-payable.wasm --is-create --input liu,wjw,123456
```

Get contract address:
```
tx GetContractAddressByTxHash -p [txHash]
tx GetContractAddressByTxHash -p 0xb57c391ee4993a1b05712806eff7646c014e29882a2062fc29249d5339a72863
```

Estimate gas:
```
tx EstimateGas -p [from],[value],[gasPrice],[gasLimit] --abi [abiPath] --wasm [wasmPath] --is-create --input [params]
tx EstimateGas -p 0x0000661A3c6c0955B5E6dbf935f0891aAA1112b9E9ca,0dip,1wu,5000000 --abi /home/qydev/testData/token-payable/token.cpp.abi.json --wasm /home/qydev/testData/token-payable/token.wasm --is-create --input liu,wjw,123456
```

Call contract:
```
tx SendTransactionContract -p [from],[contract_address],[value],[gasPrice],[gasLimit] -func-name [function_name] --input [params]
tx SendTransactionContract -p 0x0000661A3c6c0955B5E6dbf935f0891aAA1112b9E9ca,0x0014ab28B203Fd254ac6f123cC94D7a91011eFFeaf24,10dip,1wu,5000000 -func-name transfer --input 0x00005586B883Ec6dd4f8c26063E18eb4Bd228e59c3E9,100
```

Call contract without state change:
```
tx CallContract -p [from],[contract_address] -func-name [function_name] -input [params]
tx CallContract -p 0x0000661A3c6c0955B5E6dbf935f0891aAA1112b9E9ca,0x0014ab28B203Fd254ac6f123cC94D7a91011eFFeaf24 -func-name getBalance -input 0x00005586B883Ec6dd4f8c26063E18eb4Bd228e59c3E9
```

Get transaction:
```
tx Transaction [txHash]
tx Transaction -p 0xf8dd21db65b2adcb5e3ed3c61475eb66a1653d309b1a82354959fdf58852f023
```

### Chain methods

Get current block:
```
chain CurrentBlock
```

Get genesis block:
```
chain GetGenesis
```

Get transaction actual fee
```
chain GetTxActualFee -p [txHash]
chain GetTxActualFee -p 0xf8dd21db65b2adcb5e3ed3c61475eb66a1653d309b1a82354959fdf58852f023
```

Suggest gas price
```
chain SuggestGasPrice
```

Get block by number:
```
chain GetBlockByNumber -p [blockNumber]
chain GetBlockByNumber -p 1
```

Get slot by number:
```
chain GetSlotByNumber -p [blockNumber]
chain GetSlotByNumber -p 1
```

Get block by block hash:
```
chain GetBlockByHash -p [blockHash]
chain GetBlockByHash -p  0x0f7057ff3e3048ed38c0ac2353e001dad6aded5d825d43fcc924a39221713e4c
```

Get receipt by tx hash
```
chain GetReceiptByTxHash -p [txHash]
chain GetReceiptByTxHash -p 0xb57c391ee4993a1b05712806eff7646c014e29882a2062fc29249d5339a72863
```

Get receipts by block number
```
chain GetReceiptsByBlockNum -p [blockNum]
chain GetReceiptsByBlockNum -p 100
```

Search logs
```
chain GetLogs -p [jsonFile]
chain GetLogs -p {"from_block":10,"to_block":10000,"addresses":["0x0010Cb4174726E90E3ce09360B5F0488Ab29Fa5aB130"],"topics":[["Transfer"]]}
chain GetLogs -p {"block_hash":"0x000023e18421a0abfceea172867b9b4a3bcf593edd0b504554bb7d1cf5f5e7b7","addresses":["0x0010Cb4174726E90E3ce09360B5F0488Ab29Fa5aB130"],"topics":[["Transfer"]]}
```

### Verifier methods

GetVerifiers:
```
verifier GetCurVerifiers
verifier GetNextVerifiers
```

GetVerifiersBySlot
```
verifier GetVerifiersBySlot -p [slotNum]
verifier GetVerifiersBySlot -p 10
```

VerifierStatus
```
verifier VerifierStatus
```

Verifier difference between two blocks
```
verifier GetBlockDiffVerifierInfo -p [blockNum]
verifier GetBlockDiffVerifierInfo -p 10
```


### Personal methods

Look up local wallet:
```
personal ListWallet
```

Look up local wallet account:

If the wallet type and path are not specified, the default wallet is displayed
```
personal ListWalletAccount -p [walletType],[walletPath]
personal ListWalletAccount -p SoftWallet,/home/qydev/tmp/dipperin_apps/default_v0/CSWallet
```

Create new wallet:
```
personal EstablishWallet -p [walletType],[walletPath],[password]
personal EstablishWallet -p SoftWallet,/tmp/TestWallet,123

```

Recovery wallet:
```
personal RestoreWallet -p [walletType],[walletPath],[password],[passpharse],[mnemonic],...,[mnemonic]
personal RestoreWallet -p SoftWallet,/tmp/TestWallet2,123,,plastic,balcony,trophy,fuel,vacant,inmate,profit,rival,mimic,cute,hurdle,pig,column,pudding,visit,edge,rhythm,armed,cook,federal,amount,stock,damp,bring
```

Open wallet:

If the wallet type and path are not specified, the default wallet is displayed
```
personal OpenWallet -p [walletType],[walletPath],[password]
personal OpenWallet -p SoftWallet,/tmp/TestWallet3,123
```

Close wallet:

If the wallet type and path are not specified, the default wallet is displayed
```
personal CloseWallet -p [walletType],[walletPath]
personal CloseWallet -p SoftWallet,/tmp/TestWallet3
```

Add account:

If the wallet type and path are not specified, the default wallet is displayed
```
personal AddAccount -p [walletType],[walletPath]
personal AddAccount -p SoftWallet,/tmp/TestWallet3
```

Get account current balance:
```
personal CurrentBalance -p [address]
personal CurrentBalance -p 0x0000e447B8B7851D3FBD5C6A03625D288cfE9Bb5eF0E
```

Get account deposit:
```
personal CurrentStake -p [address]
personal CurrentStake -p 0x0000e447B8B7851D3FBD5C6A03625D288cfE9Bb5eF0E

```

Get account reputation:
```
personal CurrentReputation -p [address]
personal CurrentReputation -p 0x0000e447B8B7851D3FBD5C6A03625D288cfE9Bb5eF0E

```

Get account nonce:
```
personal GetTransactionNonce -p [address]
personal GetTransactionNonce -p 0x00001c2beC8E0E4caac668cD75d520E41f827092Ce79
```

Get wallet nonce:
```
personal GetAddressNonceFromWallet -p [address]
personal GetAddressNonceFromWallet -p 0x00001c2beC8E0E4caac668cD75d520E41f827092Ce79
```

Set Signer
```
personal SetBftSigner -p [address]
personal SetBftSigner -p 0x00001c2beC8E0E4caac668cD75d520E41f827092Ce79
```

### Miner methods

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
miner SetMineCoinBase -p [address]
miner SetMineCoinBase -p 0x0000e447B8B7851D3FBD5C6A03625D288cfE9Bb5eF0E
```

Set miner config:
```
miner SetMineGasConfig -p [gasFloor][gasCeil]
miner SetMineGasConfig -p 100,5000000
```