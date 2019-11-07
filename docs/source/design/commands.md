# Command Line Tool
- [How to use command line tool](#How-to-use-command-line-tool)
- [How to operate Test Node](#How-to-operate-Test-Node)
- [Related Functional Operations](#Related-Functional-Operations)

## How to use command line tool

### Connect to test environment

dipperin command line tools are located in the $GOBIN directory: `~/go/bin/dipperincli`.

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
[ModuleName] [MethodName] -p [parameters]

```

### Response Info
TxInfo
```
        [the tx info is:]
       
       	TXID:	  0x9bf08cec5cdb6db2541a4e20b52e263d5be54de5ee65b1f77d5bb8b87927de9b
       	Type:     ContractCall
       	From:     0x00005e9abe7fe3ac453e187d3df6c8d9c6f106a8b024
       	To:       0x0014d5c05b6c715e86e783d7023c06cb1ab65d6ae568
       	Nonce:    5
       	GasPrice: 1
       	GasLimit: 500000
       	Hashlock: <nil>
       	Timelock: 0x0
       	Value:    1000 CSC
       	Data:     0xf841887472616e73666572ae307830303030343137394435376534354362336235344436464145463639653734366266323430453238373937388800000000000003e8
       	V:        0xcb9
       	R:        0x98388025265dc71ed37ed2a187e651f1cc4e2ee42ba0794c4ed614f2f0f6a63a
       	S:        0x4ccfa1f02f1f1d73beb9e00e1ab5a4bc0b2f89f6453cb9405f96584b28e8737d
       	HashKey:  0x    
       
       the BlockHash is:0x0000102e0c14a7089897c84de472fa041e8895bbdc43db08706991285d5227d8
       the BlockNumber is:4711
       the TxIndex is:0

```

BlockInfo
```
       Header(0x000014cec2ed5152145f06d46cc4bd19116ccafa1feba445d3c32b014594ee54):
      [	Version:	        0
      	Number:	            6381
      	Seed:				0xacbe0a45428368c13ba779feefca47dc6be9d324c1ba7b2d2d75ad8bd52fd5e6
      	PreHash:	        0x000000ea4bd207c4e5c0d0aed8ebdf8afcb240023762c1dbcc4d18c53941e093
      	Difficulty:	        0x1e159984
      	TimeStamp:	        1572867133428162000
      	CoinBase:           0x00009396c8ff89d0D77ED1A109a1bE408123b297571a
      	GasLimit        	3360000000
      	GasUsed             0
      	Nonce:		        0x0000000000000000000000000000000000000000000000000000000000007c5f
      	Bloom：         		0x0000000000000000000000000000000000000000000000000000000000000000
      	TransactionRoot:    0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421
      	StateRoot:	        0x5199a33e84053d5ecc3e859e01439dc4943dcbe31ae52bd9ab24bcf0a14a748b
      	VerificationRoot:   0xdb7b3960a252178129fece17d52208a0f22f74fd6215659dd6815c46b5eeb4ca
      	InterlinkRoot:      0x42a72fb5ebf0538be484cdb3cf7cf077edbfd0557d0818c53422f8f751391337
      	RegisterRoot     	0x939b2b49ca6e5033e115389eb2b88bd21a1b2571f44f94ac65058d2d00a387b5
      	ReceiptHash      	0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421]
      [the current Block txs is]:
      [the current Block commit address is]:
      [commit address]:0x00004179D57e45Cb3b54D6FAEF69e746bf240E287978
      [commit address]:0x00006fC7E9B39d6C00A767AAdA3e05AEA7ba8d71ED6D
      [commit address]:0x00006532255660D9e228D997dcD827DeC685b9a17ca1
```


GenesisBlockInfo
```
        Header(0x934ad33489c3ad7b465f7d9a1bb7cdaceae3623afa63543cf74e8b67f0f44a65):
       [	Version:	        0
       	Number:	            0
       	Seed:				0x0000000000000000000000000000000000000000000000000000000000000000
       	PreHash:	        0x0000000000000000000000000000000000000000000000000000000000000000
       	Difficulty:	        0x1e566611
       	TimeStamp:	        1533715688000000000
       	CoinBase:           0x00000000000000000000000000000000000000000000
       	GasLimit        	3360000000
       	GasUsed             0
       	Nonce:		        0x0000000000000000000000000000000000000000000000000000000000000000
       	Bloom：         		0x0000000000000000000000000000000000000000000000000000000000000000
       	TransactionRoot:    0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421
       	StateRoot:	        0xeffb119b9dc9b1b899ee0fde95b9d0a86bd3cef26719c66862b55a95175d410b
       	VerificationRoot:   0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421
       	InterlinkRoot:      0x0000000000000000000000000000000000000000000000000000000000000000
       	RegisterRoot     	0x19da3f2bf161a3a674df2eeec596b4845806e0b411745a99b29fda56554bec8d
       	ReceiptHash      	0x0000000000000000000000000000000000000000000000000000000000000000]
       [the current Block txs is]:
       [the current Block commit address is]:

```

ReceiptInfo
```
PostState:			0xbc7d22661db345336ea688de2d2b2a500ea242915b87d6d096fc2a3d3ad47258
       	Status: 			Successful
       	CumulativeGasUsed:	81808
       	Bloom:              0x40000000000000000000000010000000000000000000000000200000000000000000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000800000000000000000000000000000000004020000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000
       	Logs:				[
       		ContractAddress:	0x0014D5C05b6c715e86E783d7023C06CB1AB65D6Ae568
       		TopicsHash:			[0x7134692b230b9e1ffa39098904722134159652b09c5bc41d88d6698779d228ff]
       		TopicName:			Approval
       		Data:				0x00005E9abE7FE3aC453e187D3Df6C8d9c6f106A8B024,0x00004179D57e45Cb3b54D6FAEF69e746bf240E287978,1500,
       		BlockNumber:		4711
       		TxHash:				0x9bf08cec5cdb6db2541a4e20b52e263d5be54de5ee65b1f77d5bb8b87927de9b
       		TxIndex:			0
       		BlockHash:			0x0000102e0c14a7089897c84de472fa041e8895bbdc43db08706991285d5227d8
       		Index:				0
       		Removed:			false
       		 
       		ContractAddress:	0x0014D5C05b6c715e86E783d7023C06CB1AB65D6Ae568
       		TopicsHash:			[0x0b5d2220daf8f0dfd95983d2ce625affbb7183c991271f49d818b4a64a268dbb]
       		TopicName:			Tranfer
       		Data:				0x00005E9abE7FE3aC453e187D3Df6C8d9c6f106A8B024,0x00004179D57e45Cb3b54D6FAEF69e746bf240E287978,1000,
       		BlockNumber:		4711
       		TxHash:				0x9bf08cec5cdb6db2541a4e20b52e263d5be54de5ee65b1f77d5bb8b87927de9b
       		TxIndex:			0
       		BlockHash:			0x0000102e0c14a7089897c84de472fa041e8895bbdc43db08706991285d5227d8
       		Index:				1
       		Removed:			false
       		]
       	TxHash:				0x9bf08cec5cdb6db2541a4e20b52e263d5be54de5ee65b1f77d5bb8b87927de9b
       	ContractAddress:	0x0014D5C05b6c715e86E783d7023C06CB1AB65D6Ae568
       	GasUsed:			81808
       	BlockHash     		0x0000102e0c14a7089897c84de472fa041e8895bbdc43db08706991285d5227d8   
       	BlockNumber     	4711 
       	TransactionIndex 	0
```

LogInfo
```
		ContractAddress:	0x0014D5C05b6c715e86E783d7023C06CB1AB65D6Ae568
		TopicsHash:			[0x0b5d2220daf8f0dfd95983d2ce625affbb7183c991271f49d818b4a64a268dbb]
		TopicName:			Tranfer
		Data:				,0x00005E9abE7FE3aC453e187D3Df6C8d9c6f106A8B024,10000000,
		BlockNumber:		4684
		TxHash:				0x1d8728cf1f0231a1e38bf90089a9c3033301b49f5b2ec93c8aecee1ec3d51f4b
		TxIndex:			0
		BlockHash:			0x0000060af2787c151c0f3d9fab878f4c05064e1c3810dbc45b79b50be9829db0
		Index:				0
		Removed:			false
```
### Transaction methods



AnnounceERC20:
```
tx AnnounceERC20 -p [owner_address],[token_name],[token_symbol],[token_total_supply],[decimal],[gasPrice],[gasLimit]
tx AnnounceERC20 -p 0x0000661A3c6c0955B5E6dbf935f0891aAA1112b9E9ca,wjw,dip,10000,3,10wu,100000

resp:  
       txId=0x778a9ae869a1fd598743bc3c115fcd5fa820940b9bd4b0f5d8f3ade08fae3c9e
       contract NO=0x001012831027dC1d5bBce129f58c5cfa47A7d1DC1C63 
```

ERC20Transfer:
```
tx ERC20Transfer -p [contract_address],[owner],[to_address],[amount],[gasPrice],[gasLimit]
tx ERC20Transfer -p 0x0010Cb4174726E90E3ce09360B5F0488Ab29Fa5aB130,0x0000661A3c6c0955B5E6dbf935f0891aAA1112b9E9ca,0x0000970e8128aB834E8EAC17aB8E3812f010678CF791,1000,10wu,100000

resp:
       txId=0x778a9ae869a1fd598743bc3c115fcd5fa820940b9bd4b0f5d8f3ade08fae3c9e 
```

ERC20TransferFrom:
```
tx ERC20TransferFrom -p [contract_address],[owner],[from_address],[to_address],[amount],[gasPrice],[gasLimit]
tx ERC20TransferFrom -p 0x0010Cb4174726E90E3ce09360B5F0488Ab29Fa5aB130,0x0000661A3c6c0955B5E6dbf935f0891aAA1112b9E9ca,0x0000661A3c6c0955B5E6dbf935f0891aAA1112b9E9ca,0x0000970e8128aB834E8EAC17aB8E3812f010678CF791,1,10wu,100000

resp:
       txId=0x778a9ae869a1fd598743bc3c115fcd5fa820940b9bd4b0f5d8f3ade08fae3c9e 
```

ERC20Allowance:
```
tx ERC20Allowance -p [contract_address],[owner],[spender]
tx ERC20Allowance -p 0x0010Cb4174726E90E3ce09360B5F0488Ab29Fa5aB130,0x0000661A3c6c0955B5E6dbf935f0891aAA1112b9E9ca,0x0000970e8128aB834E8EAC17aB8E3812f010678CF791

resp:
       token_allowance=1800000
```

ERC20Approve:
```
tx ERC20Approve -p [contract_address],[owner],[to_address],[amount],[gasPrice],[gasLimit]
tx ERC20Approve -p 0x0010Cb4174726E90E3ce09360B5F0488Ab29Fa5aB130,0x0000661A3c6c0955B5E6dbf935f0891aAA1112b9E9ca,0x0000970e8128aB834E8EAC17aB8E3812f010678CF791,1000,10wu,100000

resp:
       txId=0x778a9ae869a1fd598743bc3c115fcd5fa820940b9bd4b0f5d8f3ade08fae3c9e 
```

ERC20Balance:
```
tx ERC20Balance -p [contract_address],[owner_address]
tx ERC20Balance -p 0x0010Cb4174726E90E3ce09360B5F0488Ab29Fa5aB130,0x0000970e8128aB834E8EAC17aB8E3812f010678CF791

resp:
       address=0x00005E9abE7FE3aC453e187D3Df6C8d9c6f106A8B024   "token balance"=180000
```

ERC20GetInfo:
```
tx ERC20GetInfo -p [contract_address]
tx ERC20GetInfo -p 0x0010Cb4174726E90E3ce09360B5F0488Ab29Fa5aB130

resp:  
       owner=0x0000661A3c6c0955B5E6dbf935f0891aAA1112b9E9ca
       name=wjw
       symbol=dip
       decimal=3
       total supply=100000
```

Register verifier:
```
tx SendRegisterTx -p [stake],[gasPrice],[gasLimit]
tx SendRegisterTx -p 1000dip,1wu,21000


resp:
       txId=0x778a9ae869a1fd598743bc3c115fcd5fa820940b9bd4b0f5d8f3ade08fae3c9e 
```

Unregister verifier:
```
tx SendCancelTx -p [gasPrice],[gasLimit]
tx SendCancelTx -p 1wu,21000

resp:
       txId=0x778a9ae869a1fd598743bc3c115fcd5fa820940b9bd4b0f5d8f3ade08fae3c9e 
```

Redemption of the deposit:
```
tx SendUnStakeTx -p [gasPrice],[gasLimit]
tx SendUnStakeTx -p 1wu,21000

resp:
       txId=0x778a9ae869a1fd598743bc3c115fcd5fa820940b9bd4b0f5d8f3ade08fae3c9e 
```

Send transaction:
```
tx SendTx -p [to],[value],[gasPrice],[gasLimit]
tx SendTx -p 0x0000970e8128aB834E8EAC17aB8E3812f010678CF791,100dip,1wu,21000

resp:
       txId=0x778a9ae869a1fd598743bc3c115fcd5fa820940b9bd4b0f5d8f3ade08fae3c9e 
```

Create contract:
```
tx SendTransactionContract -p [from],[value],[gasPrice],[gasLimit] --abi [abiPath] --wasm [wasmPath] --is-create --input [params]
tx SendTransactionContract -p 0x0000661A3c6c0955B5E6dbf935f0891aAA1112b9E9ca,0dip,1wu,5000000 --abi /home/qydev/testData/token-payable/token-payable.cpp.abi.json --wasm /home/qydev/testData/token-payable/token-payable.wasm --is-create --input liu,wjw,123456

resp:
       txId=0x778a9ae869a1fd598743bc3c115fcd5fa820940b9bd4b0f5d8f3ade08fae3c9e 
```

Get contract address:
```
tx GetContractAddressByTxHash -p [txHash]
tx GetContractAddressByTxHash -p 0xb57c391ee4993a1b05712806eff7646c014e29882a2062fc29249d5339a72863

resp:
       Contract Address=0x0014D5C05b6c715e86E783d7023C06CB1AB65D6Ae568
```

Estimate gas:
```
tx EstimateGas -p [from],[value],[gasPrice],[gasLimit] --abi [abiPath] --wasm [wasmPath] --is-create --input [params]
tx EstimateGas -p 0x0000661A3c6c0955B5E6dbf935f0891aAA1112b9E9ca,0dip,1wu,5000000 --abi /home/qydev/testData/token-payable/token.cpp.abi.json --wasm /home/qydev/testData/token-payable/token.wasm --is-create --input liu,wjw,123456

resp:
       estimated gas=10000000
```

Call contract:
```
tx SendTransactionContract -p [from],[contract_address],[value],[gasPrice],[gasLimit] -func-name [function_name] --input [params]
tx SendTransactionContract -p 0x0000661A3c6c0955B5E6dbf935f0891aAA1112b9E9ca,0x0014ab28B203Fd254ac6f123cC94D7a91011eFFeaf24,10dip,1wu,5000000 -func-name transfer --input 0x00005586B883Ec6dd4f8c26063E18eb4Bd228e59c3E9,100

resp:
       txId=0x778a9ae869a1fd598743bc3c115fcd5fa820940b9bd4b0f5d8f3ade08fae3c9e 
```

Call contract without state change:
```
tx CallContract -p [from],[contract_address] -func-name [function_name] -input [params]
tx CallContract -p 0x0000661A3c6c0955B5E6dbf935f0891aAA1112b9E9ca,0x0014ab28B203Fd254ac6f123cC94D7a91011eFFeaf24 -func-name getBalance -input 0x00005586B883Ec6dd4f8c26063E18eb4Bd228e59c3E9

resp:
       resp=0x778a9ae869a1fd598743b
```

Get transaction:
```
tx Transaction [txHash]
tx Transaction -p 0xf8dd21db65b2adcb5e3ed3c61475eb66a1653d309b1a82354959fdf58852f023

resp:
        ${TxInfo}       
```

### Chain methods

Get current block:
```
chain CurrentBlock

resp:

      Block info:
      ${BlockInfo}
```

Get genesis block:
```
chain GetGenesis

resp:
       Block info:
       ${GenesisBlockInfo}
```

Get transaction actual fee
```
chain GetTxActualFee -p [txHash]
chain GetTxActualFee -p 0xf8dd21db65b2adcb5e3ed3c61475eb66a1653d309b1a82354959fdf58852f023

resp:

       txActualFee=10000WU
```

Suggest gas price
```
chain SuggestGasPrice

resp:
       gasPrice=10000WU
```

Get block by number:
```
chain GetBlockByNumber -p [blockNumber]
chain GetBlockByNumber -p 1

resp:  
        Block info:
        ${BlockInfo}
```

Get slot by number:
```
chain GetSlotByNumber -p [blockNumber]
chain GetSlotByNumber -p 1

resp:
       slot=2
```

Get block by block hash:
```
chain GetBlockByHash -p [blockHash]
chain GetBlockByHash -p  0x0f7057ff3e3048ed38c0ac2353e001dad6aded5d825d43fcc924a39221713e4c

resp:   
        Block info:
        {BlockInfo}
```

Get receipt by tx hash
```
chain GetReceiptByTxHash -p [txHash]
chain GetReceiptByTxHash -p 0xb57c391ee4993a1b05712806eff7646c014e29882a2062fc29249d5339a72863

resp:
        ReceiptInfo:
        $ReceiptInfo}
```

Get receipts by block number
```
chain GetReceiptsByBlockNum -p [blockNum]
chain GetReceiptsByBlockNum -p 100

resp:
       ReceiptInfos:
       [${ReceitpInfo}...]
       
```

Search logs
```
chain GetLogs -p [jsonFile]
chain GetLogs -p {"from_block":10,"to_block":10000,"addresses":["0x0014D5C05b6c715e86E783d7023C06CB1AB65D6Ae568"],"topics":[["Transfer"]]}
chain GetLogs -p {"block_hash":"0x000023e18421a0abfceea172867b9b4a3bcf593edd0b504554bb7d1cf5f5e7b7","addresses":["0x0010Cb4174726E90E3ce09360B5F0488Ab29Fa5aB130"],"topics":[["Transfer"]]}

resp:
       found logs:
       [${LogInfo}...]
       
       or 
       
       logs not found
```

### Verifier methods

GetVerifiers:
```
verifier GetCurVerifiers

resp:
       Current Verifiers:
         address: 0x00006fC7E9B39d6C00A767AAdA3e05AEA7ba8d71ED6D, is_default: true
         address: 0x00006532255660D9e228D997dcD827DeC685b9a17ca1, is_default: true
         ...
         
verifier GetNextVerifiers

resp:
       Next Verifiers:
         address: 0x00006fC7E9B39d6C00A767AAdA3e05AEA7ba8d71ED6D, is_default: true
         address: 0x00006532255660D9e228D997dcD827DeC685b9a17ca1, is_default: true
         ...    

```

GetVerifiersBySlot
```
verifier GetVerifiersBySlot -p [slotNum]
verifier GetVerifiersBySlot -p 10

resp:
        GetVerifiersBySlot result:
          address: 0x00006fC7E9B39d6C00A767AAdA3e05AEA7ba8d71ED6D, is_default: true
          address: 0x00006532255660D9e228D997dcD827DeC685b9a17ca1, is_default: true
          ...  
```

VerifierStatus
```
verifier VerifierStatus

resp:
        status="Not Registered" balance=24742.79999999999102493DIP
        
        or
        
        status="Unstaked" balance=24742.79999999999102493DIP
        
        or 

        status="Registered" balance=24742.79999999999102493DIP stake=100DIP  reputation=80 is current verifier=false
        
        or
        
        status="Canceled" balance=24742.79999999999102493DIP stake=100DIP  reputation=80
```

Verifier difference between two blocks
```
verifier GetBlockDiffVerifierInfo -p [blockNum]
verifier GetBlockDiffVerifierInfo -p 10

resp:
        the MasterVerifier address is:
       	    address: 0x00005ECCF0AAa6E8F451078448a182970e80cbDd253b
        the CommitVerifier address is:
       	    address: 0x00004179D57e45Cb3b54D6FAEF69e746bf240E287978
       	    address: 0x00006fC7E9B39d6C00A767AAdA3e05AEA7ba8d71ED6D
       	    address: 0x00006532255660D9e228D997dcD827DeC685b9a17ca1
        the NotCommitVerifier address is:
       	    address: 0x00005ECCF0AAa6E8F451078448a182970e80cbDd253b
```


### Personal methods

Look up local wallet:
```
personal ListWallet

resp:
         Call ListWallet Result:
            Wallet Info: 
             WalletType:   0
        	 Path:         /Users/konggan/tmp/dipperin_apps/nodenew
             WalletName:   CSWallet
            ... 
```

Look up local wallet account:

If the wallet type and path are not specified, the default wallet is displayed
```
personal ListWalletAccount -p [walletType],[walletPath]
personal ListWalletAccount -p SoftWallet,/home/qydev/tmp/dipperin_apps/default_v0/CSWallet

resp:
        Call ListWalletAccount result: 
        	address: 0x00005E9abE7FE3aC453e187D3Df6C8d9c6f106A8B024
        	address: 0x0000D8d017eF1Bd08897E3C25A5E530f745Bd2975C87
        	address: 0x0000D73dBB184feA834032c4fACA35b019448F156b34
        	...
```

Create new wallet:
```
personal EstablishWallet -p [walletType],[walletPath],[password]
personal EstablishWallet -p SoftWallet,/tmp/TestWallet,123

resp:
       mnemonic=upset,nurse,absent,ski,ticket,crime,sister,language,inject,wave,depend,fix,menu,boy,shy,lake,honey,fuel,thumb,hen,grab,laugh,rocket,divide
```

Recovery wallet:
```
personal RestoreWallet -p [walletType],[walletPath],[password],[passpharse],[mnemonic],...,[mnemonic]
personal RestoreWallet -p SoftWallet,/tmp/TestWallet2,123,,plastic,balcony,trophy,fuel,vacant,inmate,profit,rival,mimic,cute,hurdle,pig,column,pudding,visit,edge,rhythm,armed,cook,federal,amount,stock,damp,bring

resp:
        Call RestoreWallet success
```

Open wallet:

If the wallet type and path are not specified, the default wallet is displayed
```
personal OpenWallet -p [walletType],[walletPath],[password]
personal OpenWallet -p SoftWallet,/tmp/TestWallet3,123

resp:
        Call OpenWallet success
```

Close wallet:

If the wallet type and path are not specified, the default wallet is displayed
```
personal CloseWallet -p [walletType],[walletPath]
personal CloseWallet -p SoftWallet,/tmp/TestWallet3

resp:
        Call CloseWallet success
```

Add account:

If the wallet type and path are not specified, the default wallet is displayed
```
personal AddAccount -p [walletType],[walletPath]
personal AddAccount -p SoftWallet,/tmp/TestWallet3

resp:
        Added Account Address=0x0000D73dBB184feA834032c4fACA35b019448F156b34
```

Get account current balance:
```
personal CurrentBalance -p [address]
personal CurrentBalance -p 0x0000e447B8B7851D3FBD5C6A03625D288cfE9Bb5eF0E

resp:
        balance=24742.79999999999102493DIP
```

Get account deposit:
```
personal CurrentStake -p [address]
personal CurrentStake -p 0x0000e447B8B7851D3FBD5C6A03625D288cfE9Bb5eF0E

resp:
        stake=742.79999999999102493DIP
```

Get account reputation:
```
personal CurrentReputation -p [address]
personal CurrentReputation -p 0x0000e447B8B7851D3FBD5C6A03625D288cfE9Bb5eF0E

resp:
        reputation=70
        
        or 
        
        lookup current reputation error
```

Get account nonce:
```
personal GetTransactionNonce -p [address]
personal GetTransactionNonce -p 0x00001c2beC8E0E4caac668cD75d520E41f827092Ce79

resp:
        nonce=50
```

Get wallet nonce:
```
personal GetAddressNonceFromWallet -p [address]
personal GetAddressNonceFromWallet -p 0x00001c2beC8E0E4caac668cD75d520E41f827092Ce79

resp:
        nonce=50
```

Set Signer
```
personal SetBftSigner -p [address]
personal SetBftSigner -p 0x00001c2beC8E0E4caac668cD75d520E41f827092Ce79

resp:
        Set wallet signer（default account）succeed
```

### Miner methods

Start mining:
```
miner StartMine

resp:
       Mining Started
```

Stop mining:
```
miner StopMine

resp:
       stop mining
```

Set miner address:
```
miner SetMineCoinBase -p [address]
miner SetMineCoinBase -p 0x0000e447B8B7851D3FBD5C6A03625D288cfE9Bb5eF0E

resp:
       setting CoinBase　success
```

Set miner config:
```
miner SetMineGasConfig -p [gasFloor][gasCeil]
miner SetMineGasConfig -p 100,5000000

resp:
       setting MinerGasConfig success
```