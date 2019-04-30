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


## How to operate test node

### Start-up

```
cs_ci_ex
```

### Connect to test node

```
rpc default-c ${TestServer}:10004
```

### Perform related operations

Transfer money:
```
rpc -n m0 -m ListWallet
rpc -n m0 -m ListWalletAccount -p SoftWallet,/home/qydev/tmp/dipperin_apps/default_m0/CSWallet,CSWallet
rpc -n m0 -m SendTransaction -p [from],[to],20000,10
```

### Error

When executing the above command,it is prompted that connect refused is caused by not connecting to the node.

Generally, it is caused by changes in IP or port.

If the IP of the test server is unchanged,you can log in to test the deployment interface to see the actual port number of the node.

Connect again

Test Deployment Interface Links：
`http://${TestServer}:8889/nodes`.

```
rpc default -c ${TestServer}:[actual port number]
```

## Related Functional Operations

Separate multiple parameters by','
```
rpc -m [MethodName] -p [parameters]

```

### Block chain

Get current block:
```
rpc -m CurrentBlock
```

Get genesis block:
```
rpc -m GetGenesis
```

Get block by number
```
rpc -m GetBlockByNumber -p [blockNumber]
rpc -m GetBlockByNumber -p 1
```

Get block by block hash:
```
rpc -m GetBlockByHash -p [blockHash]
rpc -m GetBlockByHash -p  0x0f7057ff3e3048ed38c0ac2353e001dad6aded5d825d43fcc924a39221713e4c
```

### Miner master

Start mining:
```
rpc -m StartMine
```

Stop mining:
```
rpc -m StopMine
```

Set miner address:
```
rpc -m SetMineCoinBase -p [address]
rpc -m SetMineCoinBase -p 0x0000e447B8B7851D3FBD5C6A03625D288cfE9Bb5eF0E
```

Send normal transaction:
```
rpc -m SendTransaction [from],[to],[value],[transactionFee],[extradata],[nonce]
rpc -m SendTransaction -p 0x00004179D57e45Cb3b54D6FAEF69e746bf240E287978,0x00007eDe4D5D808DA8a267284b38E00ABccb42889dF2,20000,10
```

Get transaction
```
rpc -m Transaction [TxId]
rpc -m Transaction -p 0xf8dd21db65b2adcb5e3ed3c61475eb66a1653d309b1a82354959fdf58852f023
```

### Wallet

Look up local wallet:
```
rpc -m ListWallet
```

Look up local wallet account:

If the wallet type and path are not specified, the default wallet is displayed
```
rpc -m ListWalletAccount -p [walletType],[walletPath]
rpc -m ListWalletAccount -p SoftWallet,/home/qydev/tmp/dipperin_apps/default_v0/CSWallet
```

Create new wallet:
```
rpc -m EstablishWallet -p [walletType],[walletPath],[password]
rpc -m EstablishWallet -p SoftWallet,/tmp/TestWallet,123

```

Recovery wallet:
```
rpc -m RestoreWallet -p [walletType],[walletPath],[password],[passpharse],[mnemonic],...,[mnemonic]
rpc -m RestoreWallet -p SoftWallet,/tmp/TestWallet2,123,,plastic,balcony,trophy,fuel,vacant,inmate,profit,rival,mimic,cute,hurdle,pig,column,pudding,visit,edge,rhythm,armed,cook,federal,amount,stock,damp,bring
```

Open wallet:

If the wallet type and path are not specified, the default wallet is displayed
```
rpc -m OpenWallet -p [walletType],[walletPath],[password]
rpc -m OpenWallet -p SoftWallet,/tmp/TestWallet3,123
```

Close wallet:

If the wallet type and path are not specified, the default wallet is displayed
```
rpc -m CloseWallet -p [walletType],[walletPath]
rpc -m CloseWallet -p SoftWallet,/tmp/TestWallet3
```

### Account

Add account:

If the wallet type and path are not specified, the default wallet is displayed
```
rpc -m AddAccount -p [walletType],[walletPath]
rpc -m AddAccount -p SoftWallet,/tmp/TestWallet3
```

Get account current balance:
```
rpc -m CurrentBalance -p [address]
rpc -m CurrentBalance -p 0x0000e447B8B7851D3FBD5C6A03625D288cfE9Bb5eF0E
```

Get account deposit:
```
rpc -m CurrentStake -p [address]
rpc -m CurrentStake -p 0x0000e447B8B7851D3FBD5C6A03625D288cfE9Bb5eF0E

```

Get account nonce:
```
rpc -m GetAddressNonceFromWallet -p [address]
rpc -m GetAddressNonceFromWallet -p 0x00001c2beC8E0E4caac668cD75d520E41f827092Ce79
```

### Transaction

Register verifier:
```
rpc -m SendRegisterTransaction -p [from],[deposit],[transactionFee],[nonce]
rpc -m SendRegisterTransaction -p 0x00004179D57e45Cb3b54D6FAEF69e746bf240E287978,100,11
```

Unregister verifier:
```
rpc -m SendCancelTransaction -p [from],[transactionFee],[nonce]
rpc -m SendCancelTransaction -p 0x00004179D57e45Cb3b54D6FAEF69e746bf240E287978,11
```

Redemption of the deposit:
```
rpc -m SendUnStakeTransaction -p [from],[transactionFee],[nonce]
rpc -m SendUnStakeTransaction -p 0x00004179D57e45Cb3b54D6FAEF69e746bf240E287978,11
```

Get transation nonce:
```
rpc -m GetTransactionNonce -p [address]
rpc -m GetTransactionNonce -p 0x00001c2beC8E0E4caac668cD75d520E41f827092Ce79
```

### Verifiers

Get Verifiers by slot:
```
rpc -m GetVerifiersBySlot -p [round]
rpc -m GetVerifiersBySlot -p 2

```

Get current slot verifiers:
```
rpc -m GetCurVerifiers
```

Get next slot verifiers:
```
rpc -m GetNextVerifiers
```

### ERC20

Create ERC20 contract:

The contract address must be saved for successful creation
```
rpc -m AnnounceERC20 -p [owner_address], [token_name], [token_symbol], [token_total_supply], [decimal],[transactionFee]
rpc -m AnnounceERC20 -p 0x0000D07252C7A396Cc444DC0196A8b43c1A4B6c53532,chain,stack,5,3,0.00001
```

Look up ERC20 information:
```
rpc -m ERC20GetInfo -p [contract_address]
rpc -m ERC20GetInfo -p 0x00100f35adf022a8aaAbef59abB97665788CDdbA30e3
```

Look up ERC20 total supply
```
rpc -m ERC20TotalSupply -p [contract_address]
rpc -m ERC20TotalSupply -p 0x00100f35adf022a8aaAbef59abB97665788CDdbA30e3
```

Look up ERC20 token name (No need) :
```
rpc -m ERC20TokenName -p [contract_address]
rpc -m ERC20TokenName -p 0x00100f35adf022a8aaAbef59abB97665788CDdbA30e3
```

Look up ERC20 token symbol (No need) :
```
rpc -m ERC20TokenSymbol -p [contract_address]
rpc -m ERC20TokenSymbol -p 0x00100f35adf022a8aaAbef59abB97665788CDdbA30e3
```

Look up ERC20 token decimals (No need) :
```
rpc -m ERC20TokenDecimals -p [contract_address]
rpc -m ERC20TokenDecimals -p 0x00100f35adf022a8aaAbef59abB97665788CDdbA30e3
```

Look up ERC20 token quota:
```
rpc -m ERC20Allowance -p [contract_address],[from_address],[to_address]
rpc -m ERC20Allowance -p 0x00100f35adf022a8aaAbef59abB97665788CDdbA30e3,0x0000D07252C7A396Cc444DC0196A8b43c1A4B6c53532,0x0000B04985A7ccc00ab023d9bC40E241F9DF0379d8c4
```

Allocation of ERC20 contract quota:
```
rpc -m ERC20Approve -p [contract_address],[from_address],[to_address],[amount],[transactionFee]
rpc -m ERC20Approve -p 0x00100f35adf022a8aaAbef59abB97665788CDdbA30e3,0x0000D07252C7A396Cc444DC0196A8b43c1A4B6c53532,
0x0000B04985A7ccc00ab023d9bC40E241F9DF0379d8c4,4,0.00001
```

Look up ERC20 balance:
```
rpc -m ERC20Balance -p [contract_address],[owner_address]
rpc -m ERC20Balance -p 0x00100f35adf022a8aaAbef59abB97665788CDdbA30e3,0x0000D07252C7A396Cc444DC0196A8b43c1A4B6c53532
```

ERC20 contract transfer:
```
rpc -m ERC20Transfer -p [contract_address],[from_address],[to_address],[amount],[transactionFee]
rpc -m ERC20Transfer -p 0x00100f35adf022a8aaAbef59abB97665788CDdbA30e3,0x0000D07252C7A396Cc444DC0196A8b43c1A4B6c53532,
0x0000B04985A7ccc00ab023d9bC40E241F9DF0379d8c4,4,0.00001
```

ERC20 contract transfer from:
```
rpc -m ERC20TransferFrom -p [contract_address],[owner_address],[from_address],[to_address],[amount],[transactionFee]
rpc -m ERC20TransferFrom -p 0x0000B04985A7ccc00ab023d9bC40E241F9DF0379d8c4,0x0000D07252C7A396Cc444DC0196A8b43c1A4B6c5353d,
0x0000D07252C7A396Cc444DC0196A8b43c1A4B6c53532,0x00100f35adf022a8aaAbef59abB97665788CDdbA30e3,4,0.00001
```